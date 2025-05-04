Feature: Create Resource
  Support create operations for Resources.

  Rule: Must have valid authentication

    Scenario: Missing authentication
      Given a name of 10 printable ASCII characters
      And no authentication
      When I create a Resource
      Then I should fail with code=unauthenticated, msg=invalid or missing authorization

    Scenario: Invalid authentication
      Given a name of 10 printable ASCII characters
      And invalid authentication
      When I create a Resource
      Then I should fail with code=unauthenticated, msg=invalid or missing authorization

  Rule: Names must be 3-100 printable ASCII characters

    Scenario: Create resource with short name
      Given a name of 3 printable ASCII characters
      And valid authentication
      When I create a Resource
      Then I should receive the Resource

    Scenario: Create resource with long name
      Given a name of 100 printable ASCII characters
      And valid authentication
      When I create a Resource
      Then I should receive the Resource

    Scenario: Create resource with invalid name
      Given a name of 10 printable non-ASCII characters
      And valid authentication
      When I create a Resource
      Then I should fail with code=invalid_argument, msg=validation error: - resource.name: value does not match regex pattern `^[ -~]+$` [string.pattern]

    Scenario: Create resource with a too short name
     Given a name of 2 printable ASCII characters
     And valid authentication
     When I create a Resource
     Then I should fail with code=invalid_argument, msg=validation error: - resource.name: value length must be at least 3 characters [string.min_len]

    Scenario: Create resource with a too long name
     Given a name of 101 printable ASCII characters
     And valid authentication
     When I create a Resource
     Then I should fail with code=invalid_argument, msg=validation error: - resource.name: value length must be at most 100 characters [string.max_len]

  Rule: Names must not be duplicated

   Scenario: Create duplicate resource
     Given valid authentication
     And an existing resource name
     When I create a Resource
     Then I should fail with code=already_exists, msg=resource name already in use
