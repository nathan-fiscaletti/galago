# Managing Responses

Responses are used to represent the data that is sent back to the client. You can read more about the [`Response` structure](https://godoc.org/github.com/nathan-fiscaletti/galago#Response) in the official documentation.

## Overview

1. [Creating a Response](#creating-a-response)
2. [Response Headers](#response-headers)
3. [Customizing Response Serializers](#customizing-response-serializers)
4. [Downloads](#downloads)

## Creating a Response

You can create a new Response using the [`NewResponse()` function](https://godoc.org/github.com/nathan-fiscaletti/galago#NewResponse). This function takes the HTTP Status Code for the response, and a map of data that will later be serialized for the Response.

```go
response := galago.NewResponse(http.StatusOK, map[string]interface{} {
    "message": "Hello, world!",
})
```

## Response Headers

You can set a header for a Response using the [`response.SetHeader(key, val)` function](https://godoc.org/github.com/nathan-fiscaletti/galago#Response.SetHeader).

```go
response.SetHeader("MyHeader", "Value")
```

## Customizing Response Serializers

A response can use a custom Serializer to override any parent Serializer. You can set the custom Serializer for the Response using the [`response.SetSerializer(serializer)` function](https://godoc.org/github.com/nathan-fiscaletti/galago#Response.SetSerializer).

```go
response.SetSerializer(JSONSerializer())
```

For more information on Serializers, see [Managing Serialization & Deserialization](./serialization.md).

## Downloads

You can configure a Response to operate as a file download using the [`response.MakeDownload` function](https://godoc.org/github.com/nathan-fiscaletti/galago#Response.MakeDownload). This will set the `Content-Disposition` header and set the response Serializer to an instance of [`DownloadSerializer()`](https://godoc.org/github.com/nathan-fiscaletti/galago#DownloadSerializer).

You can create your data map for the response using the [`Serializer.MakeRawData()` function](https://godoc.org/github.com/nathan-fiscaletti/galago#Serializer.MakeRawData) of the `DownloadSerializer()`.

```go
response := galago.NewResponse(
                http.StatusOK, galago.DownloadSerializer().MakeRawData(
                    "This is my custom file.",
                ),
            ).MakeDownload("my_file_name.txt")
```