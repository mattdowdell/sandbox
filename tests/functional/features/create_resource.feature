Feature: Create Resource
  Support create operations for Resources.

  Rule: Names must be 3-100 printable ASCII characters

    Scenario: Create resource with short name
      Given a name of 3 printable ASCII characters
      When I create a Resource
      Then I should receive the Resource

    Scenario: Create resource with long name
      Given a name of 100 printable ASCII characters
      When I create a Resource
      Then I should receive the Resource

    Scenario: Create resource with invalid name
      Given a name of 10 printable non-ASCII characters
      When I create a Resource
      Then I should receive InvalidArgument with message: validation error: - resource.name: value does not match regex pattern `^[ -~]+$` [string.pattern]

    Scenario: Create resource with a too short name
      Given a name of 2 printable ASCII characters
      When I create a Resource
      Then I should receive InvalidArgument with message: validation error: - resource.name: value length must be at least 3 characters [string.min_len]

    Scenario: Create resource with a too long name
      Given a name of 101 printable ASCII characters
      When I create a Resource
      Then I should receive InvalidArgument with message: validation error: - resource.name: value length must be at most 100 characters [string.max_len]

  Rule: Names must not be duplicated

    Scenario: Create duplicate resource
      Given an existing resource name
      When I create a Resource
      Then I should receive AlreadyExists with message: resource name already in use
