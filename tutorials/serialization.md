# Managing Serialization & Deserialization

Serializers are used for parsing data that enters your framework via a Request, and for Serializing data before it leaves your framework in the form of a Response. You can read more about the [`Serializer` structure](https://godoc.org/github.com/nathan-fiscaletti/galago#Serializer) in the official documentation.

## Overview

1. [Serializer Precedence](#serializer-precedence)
2. [Built in Serializers](#built-in-serializers)
3. [Raw Serializers](#raw-serializers)
4. [Creating a Custom Serializer](#creating-a-custom-serializer)

## Serializer Precedence

Serializers can be applied to many different components in Galago. However, only one Serializer will ever be applied to any data. This means that the Serializers applied to certain components need to have precedence over other Serializers.

The general rules of precedence with Serializers are as follows.

- **For Response Data**

   1. A Serializer applied to a Response will always be used before any other serializer.
   2. A Serializer applied to a Route will be used if no Serializer has been set for the Response.
   3. A Serializer applied to an App will be used if no Serializer has been set for the Response.
   4. The global system serializer stored in [`galago.DefaultSerializer`](https://godoc.org/github.com/nathan-fiscaletti/galago#DefaultSerializer) will be used last.

- **For Request Data**

   1. The Serializer applied to the Route will always be used before any other Serializer.
   2. A serializer applied to an App will be used if no Serializer has been set for the Route.
   3. The global system serializer stored in [`galago.DefaultSerializer`](https://godoc.org/github.com/nathan-fiscaletti/galago#DefaultSerializer) will be used last.

## Built in Serializers

There are several Serializers built into Galago. 

- [`galago.JSONSerializer()`](https://godoc.org/github.com/nathan-fiscaletti/galago#JSONSerializer)
- [`galago.DownloadSerializer()`](https://godoc.org/github.com/nathan-fiscaletti/galago#DownloadSerializer)
- [`galago.TextSerializer()`](https://godoc.org/github.com/nathan-fiscaletti/galago#TextSerializer)

## Raw Serializers

A Raw Serializer is a Serializer that will receive raw input data and serialize a special map into raw output data. It does this using a key designated to the data.

To create a Raw Serializer, use the [`NewRawSerializer(key, contentType)` function](https://godoc.org/github.com/nathan-fiscaletti/galago#NewRawSerializer). This function takes the key at which to look for the raw data and the content type to use for the data.

```go
serializer := galago.NewRawSerializer("data", "text/plain")
```

Now, any time you apply this serializer to a Response component, it will accept a map structured as follows

```go
map[string]interface{} {
    "data": "Your raw data",
}
```

When you apply it to a Request component, it will take any raw input and create the above map using it.

You can also create this map using the [`serializer.MakeRawData(data)` function](https://godoc.org/github.com/nathan-fiscaletti/galago#Serializer.MakeRawData).

Both the `DownloadSerializer()` and the `TextSerializer()` built in Serializers are Raw Serializers.

## Creating a Custom Serializer

You can create your own custom Serializer by implementing the `Serializer` structure.

```go
serializer := &galago.Serializer{
    ContentType: "text/plain",
    Serialize: func(data map[string]interface{}) (string, error) {
        // serialize the map to a string
    },
    Deserializer: func(data string) (map[string]interface{}, error) {
        // deserialize the string to a map
    }
}
```