# [JSON and GO](https://go.dev/blog/json)

## Introduction
JSON represents objects that are used for serialization of state in backends

## Encoding
- func Marshal(v interface{}) ([]byte, error)
- encode JSON using the Marshal function
- example
    - type Message struct {
        Name string
        Body string
        Time int64
    }
    - marshal it using json.Marshal
    - parses it into a
      []byte(`"Name: testname", "Body":"Hello", "Time":1294706395881547000}`)
- Types supported
    - map[string]T
    - pointers will be encoded as the values they point to or 'null' if the pointer is nil
    - Cyclic data structures are not supported - will cause it to go into inifinte loop
    - Channel, complex and func types cannot be encoded

## Decoding
- func Unmarshal(data []byte, v interface{}) error
- parse bytes into the pointer of a struct 
- if the data is valid to fit into the struct type
    - will get parsed into the respective fields
    - identifies fields using
        - exported field with a tag `json:"Foo"`
        - exported field name "Foo"
        - an exported field with same name no case sensitivity "fOo" etc.

## Streaming Encoders and Decoders
func NewDecoder(r io.Reader) *Decoder
func NewEncoder(w io.Writer) *Encoder
- read and write streams of json data
  var v map[string]interface{}
  dec.Decode(&v)

- interface{}
    - can be used to Unmarshal arbitrary values without knowing structure
    - this works because every go type implements atleast 0 functions
        - so all types will be of type interface{} 
        - var i interface{}
        - i.(int64) to access the underlying concrete type
        SO THE MAIN POINT
        - the map[string]interface{} is used by json.Unmarshal if passed a interaface{}

        
