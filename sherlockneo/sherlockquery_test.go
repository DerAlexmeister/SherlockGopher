package sherlockneo

import (
	"testing"
)

func Test_getDropGraph(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDropGraph(); got != tt.want {
				t.Errorf("getDropGraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetContains(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getContains(); got != tt.want {
				t.Errorf("GetContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAddNode(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAddNode(); got != tt.want {
				t.Errorf("GetAddNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConstrains(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConstrains(); got != tt.want {
				t.Errorf("GetConstrains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetReturnAll(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getReturnAll(); got != tt.want {
				t.Errorf("GetReturnAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConnectbyLink(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConnectbyLink(); got != tt.want {
				t.Errorf("GetConnectbyLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCountNumberOfNodes(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCountNumberOfNodes(); got != tt.want {
				t.Errorf("GetCountNumberOfNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCountRelsToNodes(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCountRelsToNodes(); got != tt.want {
				t.Errorf("GetCountRelsToNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCountCSSNodes(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCountCSSNodes(); got != tt.want {
				t.Errorf("GetCountCSSNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCountJavascriptNodes(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCountJavascriptNodes(); got != tt.want {
				t.Errorf("GetCountJavascriptNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCountImageNodes(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCountImageNodes(); got != tt.want {
				t.Errorf("GetCountImageNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCountHtmlsNodes(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCountHtmlsNodes(); got != tt.want {
				t.Errorf("GetCountHtmlsNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResponseTimeInTableAndStatusCode(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getResponseTimeInTableAndStatusCode(); got != tt.want {
				t.Errorf("GetResponseTimeInTableAndStatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllRels(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAllRels(); got != tt.want {
				t.Errorf("GetAllRels() = %v, want %v", got, tt.want)
			}
		})
	}
}
