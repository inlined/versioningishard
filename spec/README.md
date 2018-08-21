# CloudEvents Micro - Version 1.0 #

## Abstract

CloudEvents is a vendor-neutral specification for defining the format of event data.

## Type System

The following abstract data types are available for use in attributes.

* `String` - Sequence of printable Unicode characters.

## Context Attributes

Every CloudEvent conforming to this specification MUST include one or more of the following context attributes.

These attributes, while descriptive of the event, are designed such that they can be serialized independent of the event data. This allows for them to be inspected at the destination without having to deserialize the event data.

The choice of serialization mechanism will determine how the context attributes and the event data will be materialized. For example, in the case of a JSON serialization, the context attributes and the event data might both appear within the same JSON object.

### eventID
* Type: `String`
* Description: ID of the event. The semantics of this string are explicitly undefined to ease the implementation of producers. Enables deduplication.
* Examples:
  * A database commit ID
* Constraints:
  * REQUIRED
  * MUST be a non-empty string
  * MUST be unique within the scope of the producer
  