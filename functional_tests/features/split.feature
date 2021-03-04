Feature: TrafficSplitter
  In order to test the TrafficTarget
  As a developer
  I need to ensure the specification is accepted by the server

  Scenario: Apply TrafficSplitter creates a ServiceSplitter in Consul
    Given the server is running
    And the following Consul ServiceDefaults exists
      ```
      {
          "Kind": "service-defaults",
          "Name": "myService",
          "Protocol": "http",
          "MeshGateway": {},
          "Expose": {}
      }
      ```
    And the following Consul ServiceResolver exists
      ```
      {
          "Kind": "service-resolver",
          "Name": "myService",
          "DefaultSubset": "api-primary",
          "Subsets": {
              "api-canary": {
                  "Filter": "Service.ID not contains \"api-primary\"",
                  "OnlyPassing": true
              },
              "api-primary": {
                  "Filter": "Service.ID contains \"api-primary\"",
                  "OnlyPassing": true
              }
          }
      }
      ```
    When I create a TrafficSplitter
    Then I expect a ServiceSplitter to have been created in Consul