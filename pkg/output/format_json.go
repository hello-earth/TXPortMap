package output

import (
	jsoniter "github.com/json-iterator/go"
)

// TODO

func (w *StandardWriter) formatJSON(output *ResultEvent) ([]byte, error) {
	return jsoniter.Marshal(output)
}

func (w *StandardWriter) formatSuccessJSON(output *ResultSuccess) ([]byte, error) {
	return jsoniter.Marshal(output)
}
