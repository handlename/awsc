package aws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_App_IsReadonly(t *testing.T) {
	tests := []struct {
		argv    []string
		want    bool
		wantErr bool
		errBody string
	}{
		// true
		{
			argv: []string{"ec2", "describe-instances"},
			want: true,
		},
		{
			argv: []string{"lambda", "get-function"},
			want: true,
		},
		{
			argv: []string{"iam", "list-roles"},
			want: true,
		},
		{
			argv: []string{"s3", "ls", "s3://bucket"},
			want: true,
		},
		{
			argv: []string{"s3api", "select-object-content", "s3://bucket/file.txt", "--expression", "SELECT * FROM S3Object"},
			want: true,
		},
		{
			argv: []string{"help"},
			want: true,
		},
		{
			argv: []string{"ec2", "help"},
			want: true,
		},
		{
			argv: []string{"ec2", "delete-cluster", "help"},
			want: true,
		},

		// false
		{
			argv: []string{"ec2", "delete-vpc", "--vpc-id", "vpc-12345678"},
			want: false,
		},
		{
			argv: []string{"s3", "cp", "file.txt", "s3://bucket/file.txt"},
			want: false,
		},

		// error
		{
			argv:    []string{"", ""},
			wantErr: true,
			errBody: "is empty",
		},
		{
			argv:    []string{"ec2", ""},
			wantErr: true,
			errBody: "is empty",
		},
		{
			argv:    []string{"ec2"},
			wantErr: true,
			errBody: "too short",
		},
	}

	asvc := Service{}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%t: %v", tt.want, tt.argv), func(t *testing.T) {
			got, err := asvc.IsReadonly(tt.argv)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errBody)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
