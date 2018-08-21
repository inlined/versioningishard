# Codelab: Versioning is Hard (aka the "BEF theorem")

## Background

You may have installed dependencies before that use [semver], but have you built a
software platform that is semantically versioned? What does this mean? Why do I claim
this is hard?

A group I work with, [CloudEvents], has started to see some of these issues when evaluating
how to represent CloudEvents in [Protocol Buffers]. Even I shudder sometimes at
the [proto language]. Proto generates libraries for both the binary Protocol Buffer binary format
and JSON. Those libraries and the JSON format can sometimes feel clunky, but I've slowly seen
wisdom in their ways. <!-- Proto is minimal and designed to help steer developers away from
pitfalls of cleverness. -->

![Cartoon showing two boxes labeled "needles" and "poison tipped needles".
Boxes are captioned "Clunkiness of JSON in proto-compatible APIs" and "Landmines in
proto-incompatible APIs" respectively](data/needles.jpg)

This code lab shows one of the API pitfalls that Proto guarded us against. To be more
approachable, we'll reproduce this bug without ever using or mentioning proto again. We'll
show how the goals of forwards compatibility and extensibility are different and can
even be at odds.

## Experiment

### Personas

We will wear the hats of three different personas represented by three different directories:

* Working Group ([`spec`](/spec)): The CloudEvents working group will release version 1.0 and 1.1 of the spec
* Library Vendor ([`lib`](/lib)): An OSS contributor will write a library to help devs use CloudEvents according to spec.
* App Author ([`app`](/app)): Someone who uses CloudEvents by depending on `lib/`. 

Additional data for this experiment (e.g. sample payloads) can be found in `data/`.

This repo has been annotated with tags. Each tag has a prefix of who authored the commit (e.g. `spec-1.1`). We'll use
this to see how each persona's work affects the others.

## Get started

This demo requires [Go](golang.org) and assumes you have the `GOPATH` environment variable set.
To get started, run the following commands:

```bash
mkdir -p ${GOPATH}/src/github.com/inlined
git clone git@github.com:inlined/versioningishard.git ${GOPATH}/src/github.com/inlined/versioningishard
cd ${GOPATH}/src/github.com/inlined/versioningishard
# This repository uses tags to help you walk through the steps and compare changes
git checkout app-1.0
go install github.com/inlined/versioningishard/app/versioningishard
```

Congratulations! You have the whole world running at version 1.0. The [spec 1.0] is still very
simple<sup>[1](#foot1)</sup>:
 * Any JSON object is a CloudEvent if it has an `eventId` string field
 * Extra fields may be present; these are "extensions" to the spec.

The library author has released [their support](/lib/cloudevents)

The app author uses CloudEvents 1.0 and also a "sampledRate" field because they downsample their data.
Try the whole thing together!

```bash
versioningishard --data data/event-1.0.json
versioningishard --data data/event-1.1.json
```

You should see the following output:

```bash
> versioningishard --data data/event-1.0.json
I got event 123
(It was sampled at a rate of 1 in 10)

> versioningishard --data data/event-1.1.json
I got event abc
(It was sampled at a rate of 1 in 20)
```

### Upgrades

Good news! The CloudEvents working group just released [spec 1.1]! Here's the version notes:

* Adds a new optional optional attribute `eventTime`
* Many developers have used `sampledRate` as an extension. It is now formally an optional attribute of
  the spec.

We can now upgrade our library to support the new spec:

```bash
git checkout lib-1.1
```

Now we can put on our app developer hats. Let's upgrade our dependency (it's only a minor change after all) and rerun
the program. We'll see the following output:

```bash
go install github.com/inlined/versioningishard/app/versioningishard
versioningishard --data data/event-1.0.json
versioningishard --data data/event-1.1.json
```

It looks like our output changed:

```bash
> versioningishard --data data/event-1.0.json
I got event 123
(It was not sampled)

> versioningishard --data data/event-1.1.json
I got event abc
(It was not sampled)
```

What just happened?! The code broke after a minor change! `git show lib-1.1` does not show any obvious culprit.

## Postmortem

Lets define three desirable features of a data format:

* Structured: simple (C-style) structs have many advantages. These "passive" data objects are transparent to developers,
  easily crafted in tests, terse in representation, and well suited to compiler checks or code-complete. These all
  happen by separating the definition of data from the scheme of the data (e.g. a struct def) which defines a contract
  compilers or tests can look up or enforce.
* Forwards compatible: Official maintainers can add features without breaking existing dependencies.
* Extensible: the ability for third-parties to offer features without changing the version of a product.

Just like the [CAP theorem] defines trade-offs we must make in databases (CP or AP), we must choose at most two of these
three attributes in most languages<sup>[2](#foot2)</sup>:

* Forwards compatible + Structured: This is one of the most common choices. If a library developer released a struct-
  based library, all new versions are forwards-compatible as long as new properties are not
  required<sup>[3](#fooot3)</sup>
* Forwards compatible + Extensible: Our demo _spec_ chose this route for JSON but our accidentally dropped forwards-
  compatibility by using a `struct`. Our _library_ MAY choose this route as well with some non-obvious
  work<sup>[4](#foot4)</sup>. Any library that used the most common JSON idioms in Go would break users in the next
  upgrade.
* Structured + Extensible: This is the option the library chose in this codelab. Standard features get typed support and
  extensions are possible, but extensions can only be promoted to standard features as a breaking change.

### Back to proto

So what did any of this have to do with [Protocol Buffers]? Well, protobuffers force you to choose Structured as one of
your three options. Next we can choose to either favor extensibility or forwards compatibility. The only workaround is
that in version 1.1 of the spec, a promoted extension MUST still appear in the extensions map and MAY be
included as a struct property. This workaround is impractical because there would be no motivation for users to read
from the struct field (older senders will not set it) and thus little motiviation to start setting it.

## Footnotes

<a name="foot1">1</a>: This is the simplest JSON format but not the simplest JSON spec. Golang's [encoding/json] cannot
yet support parsing unknown inline attributes ([Issue 63213]). To work around this, we use [@duglin]'s jsonext tool.
Thanks Doug! 

<a name="foot2">2</a>: Technically speaking, untyped languages make us give up structured support. Some untyped
languages have varying support of some of the Structured goal, and some struggle. For example, TypeScript has the
ability to define a general property bag with some well-known attributes, but the all values-types in the interface must
be the same type. Similarly, there's no safety net if you accidentally misspell one of the well-known attributes.

<a name="foot3">3</a>: To be totally fair, Golang has a function-style constructor in addition to the struct-style
constructor. The function-style constructor lets consumers of the struct opt-out of some types of forwards
compatibility.

<a name="foo4">4</a>: TL;DR: the library would have to be built off of a map:

```golang
type CloudEvent map[string]interface{}

func (c CloudEvent) GetEventID() string {
  return c["eventId"]
}

func (c CloudEvent) SetEventID(id string) {
  c["eventId"] = id
}

// ...
```

Now the usage `event["sampledRate"]` is always safe.


[Semver]: https://semver.org/]
[CloudEvents]: https://github.com/cloudevents/spec
[Protocol Buffers]: https://en.wikipedia.org/wiki/Protocol_Buffers
[proto language]: https://developers.google.com/protocol-buffers/docs/proto
[sampling extension]: https://github.com/cloudevents/spec/blob/master/extensions/sampled-rate.md
[encoding/json]: https://golang.org/pkg/encoding/json/
[Issue 63213]: https://github.com/golang/go/issues/6213
[@duglin]: https://github.com/duglin
[spec 1.0]: https://github.com/inlined/versioningishard/tags/spec-1.0
[spec 1.1]: https://github.com/inlined/versioningishard/tags/spec-1.1
[Thrift]: https://thrift.apache.org/
[Bond]: https://microsoft.github.io/bond/
[CAP theorem]: https://en.wikipedia.org/wiki/CAP_theorem