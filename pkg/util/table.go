/*
Copyright © 2020 Marvin

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package util

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

type Table struct {
	*tablewriter.Table
}

func NewMarkdownTableStyle(wt io.Writer, header []string, data [][]string) {
	table := tablewriter.NewWriter(wt)
	table.SetHeader(header)
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetAlignment(tablewriter.ALIGN_LEFT) // Set Alignment
	table.AppendBulk(data)
	table.Render()
}

func NewColMarkdownTable(wt io.Writer, header []string, data [][]string) {
	table := tablewriter.NewWriter(wt)
	table.SetColWidth(3000)
	table.SetHeader(header)
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: true})
	table.AppendBulk(data)
	table.Render()
}
