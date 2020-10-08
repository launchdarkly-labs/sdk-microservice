# LaunchDarkly SDK Microservice

[![CircleCI](https://circleci.com/gh/launchdarkly/sdk-agent.svg?style=svg)](https://circleci.com/gh/launchdarkly/sdk-agent)

The LaunchDarkly SDK Microservice exposes functionality of the [Go server-side SDK](https://github.com/launchdarkly/go-server-sdk) over HTTP. The SDK Microservice is intended to be run on a local network, and allows LaunchDarkly to be utilized in otherwise unsupported languages, and environments. If a native SDK exists for your platform / language, the native SDK should be used instead of SDK Microservice.

The SDK Microservice is architecturally different from the [LaunchDarkly Relay Proxy](https://github.com/launchdarkly/ld-relay). Whereas the Relay Proxy is connected to by SDKs, SDK Microservice exposes an SDK.

## Building

The SDK Microservice can be built with: `go build .` which will produce the executable `sdk-microservice`.

## Specifying a configuration

The SDK Microservice is managed with environment variables. You must set `SDK_KEY` to your LaunchDarkly projects SDK key, and may optionally set `PORT` which defaults to `8080`.

## API

## GET /

Returns HTTP 200 with a body of the form:

```
{
    "initialized": Boolean
}
```

## POST /track

Submit an event to be tracked. POST body must be of the form:

```
{
    "user": User,
    "key": String,
    "data": optional Any,
    "metricValue": optional Number
}
```

Returns HTTP 204.

## POST /flush

Flush events currently in the queue. Returns HTTP 204.

## POST /identify

Submit a user to be recorded by LaunchDarkly. POST body must be of the form:

```
{
    "user": User
}
```

## POST /allFlags

Return a map of all flags evaluated for a specific user. POST body must be of the form:

```
{
    "user": User
}
```

Returns HTTP 200 of the form:

```
{
    "flag1": Any,
    "flag2": Any
}
```

## POST /feature/{key}/eval

Evaluate a flag for a given user: Post body must be of the form:

```
{
    "user": User,
    "defaultValue": Any,
    'detail': Optional Boolean
}
```

Returns HTTP 200 of the form:

```
{
    "key": String,
    "result": Any,
    "variationIndex": Optional Number,
    "reason": Optional Reason
}
```

The fields `reason`, and `variationIndex` will only be included if you specify `detail == true` in the request.

## Type User

A `User` has the form:

```
{
    "key": String,
    "ip": Optional String,
    "firstName": Optional String,
    "lastName": Optional String,
    "email": Optional String,
    "name": Optional String,
    "avatar": Optional String,
    "country": Optional String,
    "privateAttributeNames": Optional List of Strings,
    "custom": Optional Object
}
```

## Type Reason

A `Reason` has the form:

```
{
    "kind": String
}
```

## LaunchDarkly overview

[LaunchDarkly](https://www.launchdarkly.com) is a feature management platform that serves over 100 billion feature flags daily to help teams build better software, faster. [Get started](https://docs.launchdarkly.com/docs/getting-started) using LaunchDarkly today!

[![Twitter Follow](https://img.shields.io/twitter/follow/launchdarkly.svg?style=social&label=Follow&maxAge=2592000)](https://twitter.com/intent/follow?screen_name=launchdarkly)

## Contributing

We encourage pull requests and other contributions from the community. Check out our [contributing guidelines](CONTRIBUTING.md) for instructions on how to contribute to this repository.

## About LaunchDarkly

* LaunchDarkly is a continuous delivery platform that provides feature flags as a service and allows developers to iterate quickly and safely. We allow you to easily flag your features and manage them from the LaunchDarkly dashboard.  With LaunchDarkly, you can:
    * Roll out a new feature to a subset of your users (like a group of users who opt-in to a beta tester group), gathering feedback and bug reports from real-world use cases.
    * Gradually roll out a feature to an increasing percentage of users, and track the effect that the feature has on key metrics (for instance, how likely is a user to complete a purchase if they have feature A versus feature B?).
    * Turn off a feature that you realize is causing performance problems in production, without needing to re-deploy, or even restart the application with a changed configuration file.
    * Grant access to certain features based on user attributes, like payment plan (eg: users on the ‘gold’ plan get access to more features than users in the ‘silver’ plan). Disable parts of your application to facilitate maintenance, without taking everything offline.
* LaunchDarkly provides feature flag SDKs for a wide variety of languages and technologies. Check out [our documentation](https://docs.launchdarkly.com/docs) for a complete list.
* Explore LaunchDarkly
    * [launchdarkly.com](https://www.launchdarkly.com/ "LaunchDarkly Main Website") for more information
    * [docs.launchdarkly.com](https://docs.launchdarkly.com/  "LaunchDarkly Documentation") for our documentation and SDK reference guides
    * [apidocs.launchdarkly.com](https://apidocs.launchdarkly.com/  "LaunchDarkly API Documentation") for our API documentation
    * [blog.launchdarkly.com](https://blog.launchdarkly.com/  "LaunchDarkly Blog Documentation") for the latest product updates
    * [Feature Flagging Guide](https://github.com/launchdarkly/featureflags/  "Feature Flagging Guide") for best practices and strategies
