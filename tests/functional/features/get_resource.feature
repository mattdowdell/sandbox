Feature: Get Resource
  Support get operations for Resources.

  Rule: IDs must be a UUID

    Scenario: Get resource with a valid ID
      Given the ID of an existing Resource
      When I get the Resource
      Then I should receive the Resource

    Scenario: Get non-existent resource
      Given a nil UUID
      When I get the Resource
      Then I should receive NotFound with message: TODO

    Scenario: Get resource with invalid ID
      Given an invalid ID
      When I get the Resource
      Then I should receive InvalidArgument with message: TODO
