Feature: Delete Resource
  Support delete operations for Resources.

    Scenario: Delete resource with a valid ID
      Given an existing Resource ID
      When I delete the Resource
      Then I should succeed

    Scenario: Delete non-existent resource
      Given a non-existent Resource ID
      When I delete the Resource
      Then I should fail with code=not_found, msg=resource does not exist

    Scenario: Delete resource with invalid ID
      Given an invalid Resource ID
      When I delete the Resource
      Then I should fail with code=invalid_argument, msg=validation error: - id: value must be a valid UUID [string.uuid]
