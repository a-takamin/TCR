package domain

import "testing"

func TestValidateNameSpace(t *testing.T) {
	tests := []struct {
		testName  string
		namespace string
		want      error
	}{
		{
			testName:  "スラッシュが 1 つの namespace",
			namespace: "myorg/myrepo",
			want:      nil,
		},
		// TODO: 異常系
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := ValidateNameSpace(tt.namespace)
			if got != nil {
				t.Fatalf("got is %s, but want nil", got.Error())
			}
		})
	}

}
