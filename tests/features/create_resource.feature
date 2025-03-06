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
      Given a name of 3 printable non-ASCII characters
      When I create a Resource
      Then I should receive InvalidArgument

  Rule: Names must not be duplicated

    Scenario: Create duplicate resource
      Given an existing resource name
      When I create a Resource
      Then I should receive AlreadyExists
