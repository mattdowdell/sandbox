Feature: Update Resource
  Support update operations for Resources.

  Rule: Must have valid authentication

    Scenario: Update resource without authentication
      Given no authentication
      And a non-existent Resource ID
      And a name of 10 printable ASCII characters
      When I update the Resource
      Then I should fail with code=unauthenticated, msg=invalid or missing authorization

    Scenario: Update resource with invalid authentication
      Given invalid authentication
      And a non-existent Resource ID
      And a name of 10 printable ASCII characters
      When I update the Resource
      Then I should fail with code=unauthenticated, msg=invalid or missing authorization

  Rule: Names must be 3-100 printable ASCII characters

    Scenario: Update resource with short name
      Given valid authentication
      And an existing Resource ID
      And a name of 3 printable ASCII characters
      When I update the Resource
      Then I should receive the Resource

    Scenario: Update resource with long name
      Given valid authentication
      And an existing Resource ID
      And a name of 100 printable ASCII characters
      When I update the Resource
      Then I should receive the Resource

    Scenario: Update resource with invalid name
      Given valid authentication
      And an existing Resource ID
      And a name of 10 printable non-ASCII characters
      When I update the Resource
      Then I should fail with code=invalid_argument, msg=validation error: resource.name: value does not match regex pattern `^[ -~]+$`

    Scenario: Update resource with a too short name
      Given valid authentication
      And an existing Resource ID
      And a name of 2 printable ASCII characters
      When I update the Resource
      Then I should fail with code=invalid_argument, msg=validation error: resource.name: value length must be at least 3 characters

    Scenario: Update resource with a too long name
      Given valid authentication
      And an existing Resource ID
      And a name of 101 printable ASCII characters
      When I update the Resource
      Then I should fail with code=invalid_argument, msg=validation error: resource.name: value length must be at most 100 characters

  Rule: Names must not be duplicated

    Scenario: Update duplicate resource
      Given valid authentication
      And an existing Resource ID
      And an existing resource name
      When I update the Resource
      Then I should fail with code=already_exists, msg=resource name already in use

  Rule: Only existing resources can be updated

    Scenario: Update non-existent resource
      Given valid authentication
      And a non-existent Resource ID
      And a name of 20 printable ASCII characters
      When I update the Resource
      Then I should fail with code=not_found, msg=resource does not exist

    Scenario: Update resource with invalid ID
      Given valid authentication
      And an invalid Resource ID
      And a name of 20 printable ASCII characters
      When I update the Resource
      Then I should fail with code=invalid_argument, msg=validation error: resource.id: value must be a valid UUID
