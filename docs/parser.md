## HTML Parser

### Introduction

The HTML Parser is used to transform a HTML-Document into a tree. Every level of HTML-Tags is represented as a Node in the tree. This means that a tag that is contained by another tag is a child node of the surrounding tag. Only tags are stored as Nodes. Plaintext, such as the title are stored in the title-node as its content. Attributes of tags are also stored in the Node. Thus the analyzer is able to search for a specific tag and extract its value.
***

![](https://github.com/ob-algdatii-20ss/SherlockGopher/blob/mergedwebserveranalyserandcrawler/assets/parser_sharp.png?raw=true)
***

In this shape, it is much easier for the analyzer part of SherlockGopher to scan through the document by traversing the tree. 

### Tokenizer

The tokenizer-Method serves the purpose to linearly convert the text of the HTML-Document into a list of classified tokens. Tokenclasses are 

* StartTag
* EndTag
* SelfClosingTag
* PlainText

The raw content of the tag is stored in the token for further processing.

### Parser

The Parse-Method iterates over the tokenlist. For every StartTag-Token a new Node is created, and registered as childnode of the previous node. The newly created Node is the new currentNode. An EndTag does not create a new Token, it instead sets `currentNode = currentNode.Parent()`. SelfClosing Tags do create a new Node and register it as a childnode of the currentNode but do not set the newly created Node as currentNode. A PlainText-token has the effect that its raw content is set as the content of the currentNode. 

### Malformed HTML-Documents

The parser is able to detect malformed HTML-Documents(e.g. StartTag but no EndTag or vice versa) by using a stack. For every StartTag the node is pushed onto the stack, for an EndTag the upmost node on the stack is pushed. The parser now compares the TagTypes of the pushed EndTag and the currentToken. If they match the Document is considered to be correct, if not there is a malformation. The Parser can not really correct these malformations with all certainty. It can only make a guess about how it might be intended and interpret it accordingly.
