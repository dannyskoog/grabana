package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/logs"
	"github.com/stretchr/testify/require"
)

func TestLogsPanelsCanBeDecoded(t *testing.T) {
	req := require.New(t)

	panel := DashboardLogs{
		Title:       "awesome logs",
		Description: "awesome description",
		Span:        12,
		Height:      "300px",
		Transparent: true,
		Datasource:  "some-loki",
		Repeat:      "ds",
		Targets: []LogsTarget{
			{
				Loki: &LokiTarget{
					Query:  "{namespace=\"default\"}",
					Legend: "logs",
					Ref:    "A",
				},
			},
			{
				Loki: &LokiTarget{
					Query:  "{namespace=\"other\"}",
					Legend: "other",
					Ref:    "B",
					Hidden: true,
				},
			},
		},
	}

	rowOption, err := panel.toOption()

	req.NoError(err)

	testBoard := dashboard.New("", dashboard.Row("", rowOption))
	req.Len(testBoard.Internal().Rows, 1)
	panels := testBoard.Internal().Rows[0].Panels
	req.Len(panels, 1)

	sdkPanel := panels[0]
	logsPanel := panels[0].LogsPanel

	req.NotNil(logsPanel)
	req.Len(logsPanel.Targets, 2)
	req.Equal(panel.Title, sdkPanel.Title)
	req.Equal(panel.Description, *sdkPanel.Description)
	req.Equal(panel.Span, sdkPanel.Span)
	req.True(sdkPanel.Transparent)
	req.Equal(panel.Datasource, *sdkPanel.Datasource)
	req.Equal(panel.Repeat, *sdkPanel.Repeat)
}

func TestLogsPanelsWithValidSortOrder(t *testing.T) {
	testCases := []struct {
		input    string
		expected logs.SortOrder
	}{
		{
			input:    "asc",
			expected: logs.Asc,
		},
		{
			input:    "desc",
			expected: logs.Desc,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			panel := DashboardLogs{
				Visualization: &LogsVisualization{
					Order: tc.input,
				},
			}

			rowOption, err := panel.toOption()

			req.NoError(err)

			testBoard := dashboard.New("", dashboard.Row("", rowOption))
			req.Len(testBoard.Internal().Rows, 1)
			panels := testBoard.Internal().Rows[0].Panels
			req.Len(panels, 1)

			logsPanel := panels[0].LogsPanel

			req.NotNil(logsPanel)

			req.Equal(string(tc.expected), logsPanel.Options.SortOrder)
		})
	}
}

func TestLogsPanelsWithInvalidSortOrder(t *testing.T) {
	req := require.New(t)

	panel := DashboardLogs{
		Visualization: &LogsVisualization{
			Order: "invalid",
		},
	}

	_, err := panel.toOption()
	req.Error(err)
	req.Equal(ErrInvalidSortOrder, err)
}

func TestLogsPanelsWithValidDeduplicationStrategy(t *testing.T) {
	testCases := []struct {
		input    string
		expected logs.DedupStrategy
	}{
		{
			input:    "none",
			expected: logs.None,
		},
		{
			input:    "exact",
			expected: logs.Exact,
		},
		{
			input:    "signature",
			expected: logs.Signature,
		},
		{
			input:    "numbers",
			expected: logs.Numbers,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.input, func(t *testing.T) {
			req := require.New(t)

			panel := DashboardLogs{
				Visualization: &LogsVisualization{
					Deduplication: tc.input,
				},
			}

			rowOption, err := panel.toOption()

			req.NoError(err)

			testBoard := dashboard.New("", dashboard.Row("", rowOption))
			req.Len(testBoard.Internal().Rows, 1)
			panels := testBoard.Internal().Rows[0].Panels
			req.Len(panels, 1)

			logsPanel := panels[0].LogsPanel

			req.NotNil(logsPanel)

			req.Equal(string(tc.expected), logsPanel.Options.DedupStrategy)
		})
	}
}

func TestLogsPanelsWithInvalidDedupStrategy(t *testing.T) {
	req := require.New(t)

	panel := DashboardLogs{
		Visualization: &LogsVisualization{
			Deduplication: "invalid",
		},
	}

	_, err := panel.toOption()
	req.Error(err)
	req.Equal(ErrInvalidDeduplicationStrategy, err)
}

func TestLogsPanelsWithInvalidTarget(t *testing.T) {
	req := require.New(t)

	panel := DashboardLogs{
		Targets: []LogsTarget{
			{
				// empty!
			},
		},
	}

	_, err := panel.toOption()
	req.Error(err)
	req.Equal(ErrTargetNotConfigured, err)
}

func TestLogsPanelsVisualizationOptionsCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := DashboardLogs{
		Visualization: &LogsVisualization{
			Time:           true,
			UniqueLabels:   true,
			CommonLabels:   true,
			WrapLines:      true,
			PrettifyJSON:   true,
			HideLogDetails: true,
		},
	}

	rowOption, err := panel.toOption()

	req.NoError(err)

	testBoard := dashboard.New("", dashboard.Row("", rowOption))
	req.Len(testBoard.Internal().Rows, 1)
	panels := testBoard.Internal().Rows[0].Panels
	req.Len(panels, 1)

	logsPanel := panels[0].LogsPanel

	req.False(logsPanel.Options.EnableLogDetails)
	req.True(logsPanel.Options.PrettifyLogMessage)
	req.True(logsPanel.Options.ShowLabels)
	req.True(logsPanel.Options.WrapLogMessage)
	req.True(logsPanel.Options.ShowCommonLabels)
	req.True(logsPanel.Options.ShowTime)
}
