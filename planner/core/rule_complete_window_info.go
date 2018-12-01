// Copyright 2016 PingCAP, Inc.
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

package core

import (
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/tidb/expression"
	"github.com/pingcap/tidb/types"
)

type streamWindowCompleter struct {
}

func (s *streamWindowCompleter) optimize(lp LogicalPlan) (LogicalPlan, error) {
	lp.CompleteStreamWindow()
	return lp, nil
}

func (p *LogicalProjection) CompleteStreamWindow() []*expression.Column {
	child := p.children[0]
	c := child.CompleteStreamWindow()
	if c != nil {
		p.schema.Columns = append(p.schema.Columns, c...)
		p.Exprs = append(p.Exprs, expression.Column2Exprs(c)...)
	}
	return c
}

func (la *LogicalAggregation) CompleteStreamWindow() []*expression.Column {
	//TODO: Complete here
	if la.AggWindow != nil {
		winStartCol := &expression.Column{
			ColName:  model.NewCIStr("window_start"),
			UniqueID: la.ctx.GetSessionVars().AllocPlanColumnID(),
			RetType: 	types.NewFieldType(mysql.TypeTimestamp),
		}
		winEndCol := &expression.Column{
			ColName:  model.NewCIStr("window_end"),
			UniqueID: la.ctx.GetSessionVars().AllocPlanColumnID(),
			RetType: 	types.NewFieldType(mysql.TypeTimestamp),
		}
		x := []*expression.Column{winStartCol, winEndCol}
		la.schema.Columns = append(la.schema.Columns, x...)
		return x
	}
	return nil
}

