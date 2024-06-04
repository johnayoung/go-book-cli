# Basics of YAML

YAML, which stands for "YAML Ain't Markup Language," is an easily readable data serialization format that is often compared to JSON and XML. While YAML has a more flexible and simpler syntax, which makes it more human-friendly, it is also a powerful tool for structurally representing complex data.

This section covers the fundamentals of YAML to help you quickly understand and effectively use this data serialization format. We'll focus on YAML's syntax and the key elements you'll encounter in YAML files.


## YAML Syntax

YAML's syntax is designed to be concise and easy to read. Here are the primary rules for writing YAML:

- YAML files should end in `.yaml` or `.yml`.
- A YAML file starts with `---` to indicate the start of a document within the file.
- YAML uses indentation to denote structure (similar to Python).
- Values are assigned with `:` and pairs are separated by `,`.
- YAML is case-sensitive.
- Lines that begin with `#` are comments.
- Scalar data types include strings, Booleans, and numbers.

A basic example in YAML would look like this:

```yaml
---
book: 
  title: "Pride and Prejudice"
  author: "Jane Austen"
  published_year: 1813
```

This sample YAML document represents a dictionary or object with three key-value pairs. The `---` at the start of the document is optional, but it's good practice to include it.


## Elements in YAML

YAML has three fundamental types of collection or composite data types: mapping, sequences, and scalars.

1. **Mapping**: This type is similar to what many programming languages call a dictionary, hash, or associative array, and it maps one set of values onto another set. In the example above, `book` is a mapping.

2. **Sequence**: This is a list of values that are not associated with name-value pairs as in a mapping. Here are two ways to represent a sequence (list) in YAML:

```yaml
fruits: ["apple", "banana", "cherry"]

or 

fruits: 
  - apple
  - banana
  - cherry
```
3. **Scalars**: Scalars are the basic building blocks in YAML, representing simple data types like numeric, boolean, and string data.

In the following example, a YAML document contains mappings (`book`), sequences (`chapters`), and scalar values (`"To Kill a Mockingbird"`, `"Harper Lee"`, etc.).

```yaml
---
book: 
  title: "To Kill a Mockingbird"
  author: "Harper Lee"
  published_year: 1960
  chapters: 
    - "Chapter 1"
    - "Chapter 2"
    - "Chapter 3"
```

## Conclusion

In this section, we've introduced you to the basics of YAML, focusing on its syntax and key elements. Remember that YAML uses a concise and human-friendly syntax. It features mappings, sequences, and scalars to organize data. Understanding these fundamentals will prepare you to read and write YAML files effectively.