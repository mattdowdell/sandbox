Feature: List Resources
  Support list operations for Resources.

  Rule: Must have valid authentication

    Scenario: List resources without authentication
      Given no authentication
      And a limit of 100
      When I list Resources
      Then I should fail with code=unauthenticated, msg=invalid or missing authorization

    Scenario: List resources with invalid authentication
      Given invalid authentication
      And a limit of 100
      When I list Resources
      Then I should fail with code=unauthenticated, msg=invalid or missing authorization

  Rule: Limit must be 0-100

    Scenario: List resources with a limit of 0
      Given valid authentication
      And 5 existing Resources
      And a limit of 0
      When I list Resources
      Then I should receive 0 Resources

    Scenario: List resources with a limit of 100
      Given valid authentication
      And 5 existing Resources
      And a limit of 100
      When I list Resources
      Then I should receive 5 Resources

    Scenario: List resources with a limit of 2
      Given valid authentication
      And 5 existing Resources
      And a limit of 2
      When I list Resources
      Then I should receive 2 Resources

    Scenario: List resources with none created
      Given valid authentication
      And a limit of 100
      When I list Resources
      Then I should receive 0 Resources

    Scenario: List resources with a limit of -1
      Given valid authentication
      And a limit of -1
      When I list Resources
      Then I should fail with code=invalid_argument, msg=validation error: limit: value must be greater than or equal to 0 and less than or equal to 100

    Scenario: List resources with a limit of 101
      Given valid authentication
      And a limit of 101
      When I list Resources
      Then I should fail with code=invalid_argument, msg=validation error: limit: value must be greater than or equal to 0 and less than or equal to 100

  Rule: Next token is returned when more results are available

  Rule: Sorting changes the order of results

  Rule: Filtering can remove resources from results




