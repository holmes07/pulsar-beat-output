// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package transfer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

// Please also see https://docs.aws.amazon.com/goto/WebAPI/transfer-2018-11-05/UpdateUserRequest
type UpdateUserInput struct {
	_ struct{} `type:"structure"`

	// A parameter that specifies the landing directory (folder) for a user when
	// they log in to the server using their client. An example is /home/username .
	HomeDirectory *string `type:"string"`

	// Allows you to supply a scope-down policy for your user so you can use the
	// same AWS Identity and Access Management (IAM) role across multiple users.
	// The policy scopes down user access to portions of your Amazon S3 bucket.
	// Variables you can use inside this policy include ${Transfer:UserName}, ${Transfer:HomeDirectory},
	// and ${Transfer:HomeBucket}.
	//
	// For scope-down policies, AWS Transfer for SFTP stores the policy as a JSON
	// blob, instead of the Amazon Resource Name (ARN) of the policy. You save the
	// policy as a JSON blob and pass it in the Policy argument.
	//
	// For an example of a scope-down policy, see "https://docs.aws.amazon.com/transfer/latest/userguide/users.html#users-policies-scope-down">Creating
	// a Scope-Down Policy.
	//
	// For more information, see "https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html"
	// in the AWS Security Token Service API Reference.
	Policy *string `type:"string"`

	// The IAM role that controls your user's access to your Amazon S3 bucket. The
	// policies attached to this role will determine the level of access you want
	// to provide your users when transferring files into and out of your Amazon
	// S3 bucket or buckets. The IAM role should also contain a trust relationship
	// that allows the Secure File Transfer Protocol (SFTP) server to access your
	// resources when servicing your SFTP user's transfer requests.
	Role *string `type:"string"`

	// A system-assigned unique identifier for an SFTP server instance that the
	// user account is assigned to.
	//
	// ServerId is a required field
	ServerId *string `type:"string" required:"true"`

	// A unique string that identifies a user and is associated with a server as
	// specified by the ServerId. This is the string that will be used by your user
	// when they log in to your SFTP server. This user name is a minimum of 3 and
	// a maximum of 32 characters long. The following are valid characters: a-z,
	// A-Z, 0-9, underscore, and hyphen. The user name can't start with a hyphen.
	//
	// UserName is a required field
	UserName *string `type:"string" required:"true"`
}

// String returns the string representation
func (s UpdateUserInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *UpdateUserInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "UpdateUserInput"}

	if s.ServerId == nil {
		invalidParams.Add(aws.NewErrParamRequired("ServerId"))
	}

	if s.UserName == nil {
		invalidParams.Add(aws.NewErrParamRequired("UserName"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// UpdateUserResponse returns the user name and server identifier for the request
// to update a user's properties.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/transfer-2018-11-05/UpdateUserResponse
type UpdateUserOutput struct {
	_ struct{} `type:"structure"`

	// A system-assigned unique identifier for an SFTP server instance that the
	// user account is assigned to.
	//
	// ServerId is a required field
	ServerId *string `type:"string" required:"true"`

	// The unique identifier for a user that is assigned to the SFTP server instance
	// that was specified in the request.
	//
	// UserName is a required field
	UserName *string `type:"string" required:"true"`
}

// String returns the string representation
func (s UpdateUserOutput) String() string {
	return awsutil.Prettify(s)
}

const opUpdateUser = "UpdateUser"

// UpdateUserRequest returns a request value for making API operation for
// AWS Transfer for SFTP.
//
// Assigns new properties to a user. Parameters you pass modify any or all of
// the following: the home directory, role, and policy for the UserName and
// ServerId you specify.
//
// The response returns the ServerId and the UserName for the updated user.
//
//    // Example sending a request using UpdateUserRequest.
//    req := client.UpdateUserRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/transfer-2018-11-05/UpdateUser
func (c *Client) UpdateUserRequest(input *UpdateUserInput) UpdateUserRequest {
	op := &aws.Operation{
		Name:       opUpdateUser,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &UpdateUserInput{}
	}

	req := c.newRequest(op, input, &UpdateUserOutput{})
	return UpdateUserRequest{Request: req, Input: input, Copy: c.UpdateUserRequest}
}

// UpdateUserRequest is the request type for the
// UpdateUser API operation.
type UpdateUserRequest struct {
	*aws.Request
	Input *UpdateUserInput
	Copy  func(*UpdateUserInput) UpdateUserRequest
}

// Send marshals and sends the UpdateUser API request.
func (r UpdateUserRequest) Send(ctx context.Context) (*UpdateUserResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &UpdateUserResponse{
		UpdateUserOutput: r.Request.Data.(*UpdateUserOutput),
		response:         &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// UpdateUserResponse is the response type for the
// UpdateUser API operation.
type UpdateUserResponse struct {
	*UpdateUserOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// UpdateUser request.
func (r *UpdateUserResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
