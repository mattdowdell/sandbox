Feature: Delete Resource
  Support delete operations for Resources.

  Rule: Must have valid authentication

    Scenario: Delete resource without authentication
      Given no authentication
      And a non-existent Resource ID
      When I delete the Resource
      Then I should fail with code=unauthenticated, msg=invalid or missing authorization

    Scenario: Delete resource with invalid authentication
      Given invalid authentication
      And a non-existent Resource ID
      When I delete the Resource
      Then I should fail with code=unauthenticated, msg=invalid or missing authorization

  Rule: Only existing resources can be deleted

    Scenario: Delete resource with a valid ID
      Given valid authentication
      And an existing Resource ID
      When I delete the Resource
      Then I should succeed

    Scenario: Delete non-existent resource
      Given valid authentication
      And a non-existent Resource ID
      When I delete the Resource
      Then I should fail with code=not_found, msg=resource does not exist

    Scenario: Delete resource with invalid ID
      Given valid authentication
      And an invalid Resource ID
      When I delete the Resource
      Then I should fail with code=invalid_argument, msg=validation error: id: value must be a valid UUID
