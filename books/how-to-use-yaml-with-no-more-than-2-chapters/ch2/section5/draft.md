# Implementing JSON in YAML

JSON (JavaScript Object Notation) and YAML (Yet Another Markup Language) are two popular data serialization formats.
They have similar capabilities, and YAML is often described as a superset of JSON. This means that every JSON file should also be a valid YAML file. This book section explores the implementation of JSON in YAML, delving into JSON syntax in YAML and comparing JSON and YAML data types. 

## JSON Syntax in YAML

When embedding a JSON object in a YAML document, it is crucial to understand and correctly use the syntax. Being aware of the conventions, such as quotation marks, colons, and commas can ensure a seamless translation between the two formats. 

### JSON Objects and Arrays:
The fundamental concept to grasp in YAML is the representation of JSON objects and arrays. A JSON object looks like this:

```json
{
  "name": "John",
  "age": 30,
  "city": "New York"
}
```
In YAML, the same object can be represented as:

```yaml
name: John
age: 30
city: New York
```
For JSON arrays, here's a sample JSON array:

```json
["Ford", "BMW", "Fiat"]
```
The same array can be written in this manner in YAML:

```yaml
- Ford
- BMW
- Fiat
```

### Syntax Discrepancies:
One critical part to note is the usage of quotation marks, colons, and commas. In JSON, both keys and string values are often enclosed in double quotes, while in YAML, quotes are not necessary unless the string includes reserved characters such as "#", ":", "{", "}" and so on. 

## JSON and YAML Data Types

When dealing with data types transfer between JSON and YAML, it's essential to understand that the data types in the two formats largely correspond with each other. For instance, both utilize Boolean values, null, numbers, strings, objects, and arrays. 

However, YAML offers a broader range of data types, including date, timestamp, binary, and pairs, which JSON lacks. Hence, when these unique data types are converted from YAML to JSON, they will be processed as strings rather than their original data types.

### Practical Example:

Look at this YAML example:

```yaml
boolean: true
null: null
number: 123
string: "hello"
object:
  key1: value1
  key2: value2
date: 2018-12-25
```
When translated into JSON:

```json
{
  "boolean": true,
  "null": null,
  "number": 123,
  "string": "hello",
  "object": {
    "key1": "value1",
    "key2": "value2"
  },
  "date": "2018-12-25"
}
```
Here, the date was represented as a date data type in YAML, but in JSON, it was converted to a string.

# Conclusion

To sum up, JSON can be implemented within YAML as YAML is a superset of JSON. However, to ensure a smooth implementation, understanding the representation of JSON syntax in YAML and grasping the compatibility of JSON and YAML data types is vital. By knowing these, one can succeed in leveraging the flexibility of YAML while still relying on the ubiquity and simplicity of JSON.