Feature: Delete Resource
  Support delete operations for Resources.

  Rule: IDs must be a UUID

    Scenario: Delete resource with a valid ID
      Given the ID of an existing Resource
      When I delete the Resource
      Then I should succeed

    Scenario: Delete non-existent resource
      Given a nil UUID
      When I delete the Resource
      Then I should receive NotFound with message: TODO

    Scenario: Delete resource with invalid ID
      Given an invalid ID
      When I delete the Resource
      Then I should receive InvalidArgument with message: TODO
