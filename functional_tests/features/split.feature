Feature: split.smi-spec.io
  In order to test the TrafficTarget
  As a developer
  I need to ensure the specification is accepted by the server

  Background:
    Given the following "service-defaults" config entry exists
    ```
    {
      "Kind": "service-defaults",
      "Name": "foo",
      "Protocol": "http"
    }
    ```
    And the following "service-defaults" config entry exists
    ```
    {
      "Kind": "service-defaults",
      "Name": "bar",
      "Protocol": "http"
    }
    ```
    And the following "service-defaults" config entry exists
    ```
    {
      "Kind": "service-defaults",
      "Name": "baz",
      "Protocol": "http"
    }
    ```
    And the following "service-defaults" config entry exists
    ```
    {
      "Kind": "service-defaults",
      "Name": "ab-test",
      "Protocol": "http"
    }
    ```

  @split @alpha1
  Scenario: Apply alpha1 TrafficSplit
    Given the server is running
    When I create the following resource
    ```
      apiVersion: split.smi-spec.io/v1alpha1
      kind: TrafficSplit
      metadata:
        name: trafficsplit-sample
      spec:
        service: foo
        backends:
          - service: bar
            weight: 50m
          - service: baz
            weight: 50m
    ```
    Then I expect the following "service-splitter" to be created
    ```
    {
      "Kind": "service-splitter",
      "Name": "foo",
      "Splits": [
        {
          "Weight": 50,
          "Service": "bar"
        },
        {
          "Weight": 50,
          "Service": "baz"
        }
      ]
    }
    ```
  
  @split
  Scenario: Apply alpha2 TrafficSplit
    Given the server is running
    When I create the following resource
    ```
      apiVersion: split.smi-spec.io/v1alpha2
      kind: TrafficSplit
      metadata:
        name: trafficsplit-sample
      spec:
        service: foo
        backends:
          - service: bar
            weight: 50
          - service: baz
            weight: 50
    ```
    Then I expect the following "service-splitter" to be created
    ```
    {
      "Kind": "service-splitter",
      "Name": "foo",
      "Splits": [
        {
          "Weight": 50,
          "Service": "bar"
        },
        {
          "Weight": 50,
          "Service": "baz"
        }
      ]
    }
    ```
  
  @split @alpha3
  Scenario: Apply alpha3 TrafficSplit
    Given the server is running
    When I create the following resource
    ```
      apiVersion: split.smi-spec.io/v1alpha3
      kind: TrafficSplit
      metadata:
        name: ab-test
      spec:
        service: foo
        matches:
        - kind: HTTPRouteGroup
          name: foo
          apiGroup: specs.smi-spec.io
        backends:
        - service: bar
          weight: 0
        - service: baz
          weight: 100
    ```
    And I create the following resource
    ```
      apiVersion: specs.smi-spec.io/v1alpha3
      kind: HTTPRouteGroup
      metadata:
        name: foo
      spec:
        matches:
        - name: metrics
          pathRegex: "/metrics"
          methods:
          - GET
        - name: health
          pathRegex: "/ping"
          methods: ["*"]
    ```
    Then I expect the following "service-splitter" to be created
    ```
    {
      "Kind": "service-splitter",
      "Name": "foo",
      "Splits": [
        {
          "Weight": 0,
          "Service": "bar"
        },
        {
          "Weight": 100,
          "Service": "baz"
        }
      ]
    }
    ```
    And I expect the following "service-router" to be created
    ```
    {
      "Kind": "service-router",
      "Name": "foo",
      "Routes": [
        {
          "Match": {
            "HTTP": {
              "PathRegex": "/metrics",
              "Methods": ["GET"]
            }
          },
          "Destination": {
            "Service": "foo"
          }
        },
        {
          "Match": {
            "HTTP": {
              "PathRegex": "/ping"
            }
          },
          "Destination": {
            "Service": "foo"
          }
        }
      ]
    }
    ```
  
  @split @alpha4
  Scenario: Apply alpha4 TrafficSplit
    Given the server is running
    When I create the following resource
    ```
      apiVersion: split.smi-spec.io/v1alpha4
      kind: TrafficSplit
      metadata:
        name: foo
      spec:
        service: foo
        matches:
        - kind: HTTPRouteGroup
          name: ab-test
          apiGroup: specs.smi-spec.io
        backends:
        - service: bar
          weight: 0
        - service: baz
          weight: 100
    ```
    And I create the following resource
    ```
      apiVersion: specs.smi-spec.io/v1alpha4
      kind: HTTPRouteGroup
      metadata:
        name: ab-test
      spec:
        matches:
        - name: metrics
          pathRegex: "/metrics"
          methods:
          - GET
          headers:
            x-debug: "1"
        - name: health
          pathRegex: "/ping"
          methods: ["*"]
    ```
    Then I expect the following "service-splitter" to be created
    ```
    {
      "Kind": "service-splitter",
      "Name": "foo",
      "Splits": [
        {
          "Weight": 0,
          "Service": "bar"
        },
        {
          "Weight": 100,
          "Service": "baz"
        }
      ]
    }
    ```
    And I expect the following "service-router" to be created
    ```
    {
      "Kind": "service-router",
      "Name": "ab-test",
      "Routes": [
        {
          "Match": {
            "HTTP": {
              "PathRegex": "/metrics",
              "Methods": ["GET"],
              "Headers":{
                "x-debug": "1"
              }
            }
          },
          "Destination": {
            "Service": "foo"
          }
        },
        {
          "Match": {
            "HTTP": {
              "PathRegex": "/ping"
            }
          },
          "Destination": {
            "Service": "foo"
          }
        }
      ]
    }
    ```