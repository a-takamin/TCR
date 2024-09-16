package handler

import "testing"

func TestPickUpName(t *testing.T) {
	type args struct {
		path              string
		partsNumAfterName int
	}
	type wants struct {
		name      string
		afterName string
		err       error
	}
	tests := []struct {
		testName string
		args     args
		wants    wants
	}{
		{
			testName: "name のスラッシュが 1 つで partsNumAfterName が 2 のときの正常系",
			args: args{
				path:              "org/repo/blobs/upload",
				partsNumAfterName: 2,
			},
			wants: wants{
				name:      "org/repo",
				afterName: "/blobs/upload",
				err:       nil,
			},
		},
		{
			testName: "name のスラッシュが 0 個で partsNumAfterName が 2 のときの正常系",
			args: args{
				path:              "orgrepo/blobs/upload",
				partsNumAfterName: 2,
			},
			wants: wants{
				name:      "orgrepo",
				afterName: "/blobs/upload",
				err:       nil,
			},
		},
		{
			testName: "name のスラッシュが 2 つで partsNumAfterName が 2 のときの正常系",
			args: args{
				path:              "org/repo/hoge/blobs/upload",
				partsNumAfterName: 2,
			},
			wants: wants{
				name:      "org/repo/hoge",
				afterName: "/blobs/upload",
				err:       nil,
			},
		},
		{
			testName: "name のスラッシュが 1 つで partsNumAfterName が 3 のときの正常系",
			args: args{
				path:              "org/repo/blobs/upload/uuid",
				partsNumAfterName: 3,
			},
			wants: wants{
				name:      "org/repo",
				afterName: "/blobs/upload/uuid",
				err:       nil,
			},
		},
		// TODO: 異常系
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			name, afterName, err := pickUpName(tt.args.path, tt.args.partsNumAfterName)
			if name != tt.wants.name {
				t.Fatalf("name is %s, but want %s", name, tt.wants.name)
			}
			if afterName != tt.wants.afterName {
				t.Fatalf("afterName is %s, but want %s", afterName, tt.wants.afterName)
			}
			if err != tt.wants.err {
				t.Fatalf("err is %s, but want %s", err.Error(), tt.wants.err.Error())
			}
		})
	}
}
