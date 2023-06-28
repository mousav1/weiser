package response

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	TotalRecords int    `json:"totalRecords,omitempty"`
	PageSize     int    `json:"pageSize,omitempty"`
	CurrentPage  int    `json:"currentPage,omitempty"`
	NextPageURL  string `json:"nextPageUrl,omitempty"`
}

type Response struct {
	Data               interface{} `json:"data,omitempty"`
	Message            string      `json:"message,omitempty"`
	StatusCode         int         `json:"-"`
	Headers            map[string]string
	ContentType        string
	ContentDisposition string
	Charset            string
	ProtocolVersion    string
	Pagination         *Pagination `json:"pagination,omitempty"`
}

func New(data interface{}, message string, statusCode int) *Response {
	return &Response{
		Data:       data,
		Message:    message,
		StatusCode: statusCode,
		Headers:    make(map[string]string),
	}
}

func Send(ctx *fiber.Ctx, response *Response) error {
	for key, value := range response.Headers {
		ctx.Set(key, value)
	}

	if response.ContentType != "" {
		contentType := response.ContentType
		if response.Charset != "" {
			contentType += "; charset=" + response.Charset
		}
		ctx.Set(fiber.HeaderContentType, contentType)
	}

	if response.ContentDisposition != "" {
		ctx.Set(fiber.HeaderContentDisposition, response.ContentDisposition)
	}

	if response.ProtocolVersion != "" {
		ctx.Set(fiber.HeaderServer, response.ProtocolVersion)
	}

	if response.Pagination != nil {
		ctx.Set("X-Total-Count", strconv.Itoa(response.Pagination.TotalRecords))
		ctx.Set("X-Page-Size", strconv.Itoa(response.Pagination.PageSize))
		ctx.Set("X-Current-Page", strconv.Itoa(response.Pagination.CurrentPage))
		if response.Pagination.NextPageURL != "" {
			ctx.Set("Link", fmt.Sprintf("<%s>; rel=\"next\"", response.Pagination.NextPageURL))
		}
	}

	return ctx.Status(response.StatusCode).JSON(response)
}
