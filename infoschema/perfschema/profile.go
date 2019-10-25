// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package perfschema

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/pprof/graph"
	"github.com/google/pprof/measurement"
	"github.com/google/pprof/profile"
	"github.com/google/pprof/report"
	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/sessionctx"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util/sqlexec"
)

type Node struct {
	Name      string
	Location  string
	Cum       int64
	CumFormat string
	Percent   string
	Children  []*Node
}

type profileGraphCollector struct {
	Rows [][]types.Datum
}

func (c *profileGraphCollector) collect(node *Node, depth int64, indent string, rootChild int, parentCur int64, isLastChild bool) {
	row := types.MakeDatums(
		prettyIdentifier(node.Name, indent, isLastChild),
		node.Percent,
		strings.TrimSpace(measurement.Percentage(node.Cum, parentCur)),
		rootChild,
		depth,
		node.Location,
	)
	c.Rows = append(c.Rows, row)

	indent4Child := getIndent4Child(indent, isLastChild)
	for i, child := range node.Children {
		rc := rootChild
		if depth == 0 {
			rc = i + 1
		}
		c.collect(child, depth+1, indent4Child, rc, node.Cum, i == len(node.Children)-1)
	}
}

const (
	// treeBody indicates the current operator sub-tree is not finished, still
	// has child operators to be attached on.
	treeBody = '│'
	// treeMiddleNode indicates this operator is not the last child of the
	// current sub-tree rooted by its parent.
	treeMiddleNode = '├'
	// treeLastNode indicates this operator is the last child of the current
	// sub-tree rooted by its parent.
	treeLastNode = '└'
	// treeGap is used to represent the gap between the branches of the tree.
	treeGap = ' '
	// treeNodeIdentifier is used to replace the treeGap once we need to attach
	// a node to a sub-tree.
	treeNodeIdentifier = '─'
)

func getIndent4Child(indent string, isLastChild bool) string {
	if !isLastChild {
		return string(append([]rune(indent), treeBody, treeGap))
	}

	// If the current node is the last node of the current operator tree, we
	// need to end this sub-tree by changing the closest treeBody to a treeGap.
	indentBytes := []rune(indent)
	for i := len(indentBytes) - 1; i >= 0; i-- {
		if indentBytes[i] == treeBody {
			indentBytes[i] = treeGap
			break
		}
	}

	return string(append(indentBytes, treeBody, treeGap))
}

func prettyIdentifier(id, indent string, isLastChild bool) string {
	if len(indent) == 0 {
		return id
	}

	indentBytes := []rune(indent)
	for i := len(indentBytes) - 1; i >= 0; i-- {
		if indentBytes[i] != treeBody {
			continue
		}

		// Here we attach a new node to the current sub-tree by changing
		// the closest treeBody to a:
		// 1. treeLastNode, if this operator is the last child.
		// 2. treeMiddleNode, if this operator is not the last child..
		if isLastChild {
			indentBytes[i] = treeLastNode
		} else {
			indentBytes[i] = treeMiddleNode
		}
		break
	}

	// Replace the treeGap between the treeBody and the node to a
	// treeNodeIdentifier.
	indentBytes[len(indentBytes)-1] = treeNodeIdentifier
	return string(indentBytes) + id
}

func profileReaderToDatums(f io.Reader) ([][]types.Datum, error) {
	p, err := profile.Parse(f)
	if err != nil {
		return nil, err
	}
	return profileToDatums(p)
}

func profileToDatums(p *profile.Profile) ([][]types.Datum, error) {
	err := p.Aggregate(true, true, true, true, true)
	if err != nil {
		return nil, err
	}
	rpt := report.NewDefault(p, report.Options{
		OutputFormat: report.Dot,
		CallTree:     true,
	})
	g, config := report.GetDOT(rpt)
	var nodes []*Node
	nroots := 0
	rootValue := int64(0)
	nodeArr := []string{}
	nodeMap := map[*graph.Node]*Node{}
	// Make all nodes and the map, collect the roots.
	for _, n := range g.Nodes {
		v := n.CumValue()
		node := &Node{
			Name:      n.Info.Name,
			Location:  fmt.Sprintf("%s:%d", n.Info.File, n.Info.Lineno),
			Cum:       v,
			CumFormat: config.FormatValue(v),
			Percent:   strings.TrimSpace(measurement.Percentage(v, config.Total)),
		}
		nodes = append(nodes, node)
		if len(n.In) == 0 {
			nodes[nroots], nodes[len(nodes)-1] = nodes[len(nodes)-1], nodes[nroots]
			nroots++
			rootValue += v
		}
		nodeMap[n] = node
		// Get all node names into an array.
		nodeArr = append(nodeArr, n.Info.Name)
	}
	// Populate the child links.
	for _, n := range g.Nodes {
		node := nodeMap[n]
		for child := range n.Out {
			node.Children = append(node.Children, nodeMap[child])
		}
	}

	rootNode := &Node{
		Name:      "root",
		Location:  "root",
		Cum:       rootValue,
		CumFormat: config.FormatValue(rootValue),
		Percent:   strings.TrimSpace(measurement.Percentage(rootValue, config.Total)),
		Children:  nodes[0:nroots],
	}

	c := profileGraphCollector{}
	c.collect(rootNode, 0, "", 0, config.Total, len(rootNode.Children) == 0)
	return c.Rows, nil
}

func cpuProfileGraph() ([][]types.Datum, error) {
	buffer := &bytes.Buffer{}
	if err := pprof.StartCPUProfile(buffer); err != nil {
		panic(err)
	}
	time.Sleep(20 * time.Second)
	pprof.StopCPUProfile()
	return profileReaderToDatums(buffer)
}

// TODO: use cluster info to get all tikv profile
func tikvCpuProfileGraph(ctx sessionctx.Context) ([][]types.Datum, error) {
	sql := "SELECT name, address, status_address FROM INFORMATION_SCHEMA.TIDB_CLUSTER_INFO WHERE type='tikv'"
	rows, _, err := ctx.(sqlexec.RestrictedSQLExecutor).ExecRestrictedSQL(sql)
	if err != nil {
		return nil, err
	}

	type result struct {
		rows [][]types.Datum
		err  error
	}

	var finalRows [][]types.Datum
	wg := sync.WaitGroup{}
	ch := make(chan result, len(rows))
	for _, row := range rows {
		name := row.GetString(0)
		address := row.GetString(1)
		statusAddr := row.GetString(2)
		if len(statusAddr) == 0 {
			ctx.GetSessionVars().StmtCtx.AppendWarning(errors.Errorf("tikv node %s(%s) does not contain status address", name, address))
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(fmt.Sprintf("http://%s/pprof/cpu?seconds=20", statusAddr))
			if err != nil {
				ch <- result{err: err}
				return
			}
			defer resp.Body.Close()
			rows, err := profileReaderToDatums(resp.Body)
			if err != nil {
				ch <- result{err: err}
				return
			}
			// add extra info
			for i := range rows {
				rows[i] = append(types.MakeDatums(name, address), rows[i]...)
			}
			ch <- result{rows: rows}
		}()
	}

	wg.Wait()
	close(ch)
	for result := range ch {
		if result.err != nil {
			ctx.GetSessionVars().StmtCtx.AppendWarning(err)
			continue
		}
		finalRows = append(finalRows, result.rows...)
	}
	return finalRows, nil
}

func profileGraph(name string) ([][]types.Datum, error) {
	p := pprof.Lookup(name)
	if p == nil {
		return nil, errors.Errorf("cannot retrieve %s profile", name)
	}
	buffer := &bytes.Buffer{}
	if err := p.WriteTo(buffer, 0); err != nil {
		return nil, err
	}
	return profileReaderToDatums(buffer)
}

func goroutinesList() ([][]types.Datum, error) {
	p := pprof.Lookup("goroutine")
	if p == nil {
		return nil, errors.Errorf("cannot retrieve goroutine profile")
	}

	buffer := bytes.Buffer{}
	err := p.WriteTo(&buffer, 2)
	if err != nil {
		return nil, err
	}

	goroutines := strings.Split(buffer.String(), "\n\n")
	var rows [][]types.Datum
	for _, goroutine := range goroutines {
		colIndex := strings.Index(goroutine, ":")
		if colIndex < 0 {
			return nil, errors.New("goroutine incompatible with current go version")
		}

		headers := strings.SplitN(strings.TrimSpace(goroutine[len("goroutine")+1:colIndex]), " ", 2)
		if len(headers) != 2 {
			return nil, errors.Errorf("incompatible goroutine headers: %s", goroutine[len("goroutine")+1:colIndex])
		}

		id, err := strconv.Atoi(strings.TrimSpace(headers[0]))
		if err != nil {
			return nil, errors.Annotatef(err, "invalid goroutine id: %s", headers[0])
		}
		state := strings.Trim(headers[1], "[]")
		stack := strings.Split(strings.TrimSpace(goroutine[colIndex+1:]), "\n")
		for i := 0; i < len(stack)/2; i++ {
			fn := stack[i*2]
			loc := stack[i*2+1]
			var identifier string
			if i == 0 {
				identifier = fn
			} else if i == len(stack)/2-1 {
				identifier = string(treeLastNode) + string(treeNodeIdentifier) + fn
			} else {
				identifier = string(treeMiddleNode) + string(treeNodeIdentifier) + fn
			}
			rows = append(rows, []types.Datum{
				types.NewStringDatum(identifier),
				types.NewIntDatum(int64(id)),
				types.NewStringDatum(state),
				types.NewStringDatum(strings.TrimSpace(loc)),
			})
		}
	}
	return rows, nil
}
