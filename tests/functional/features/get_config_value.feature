Feature: Get Config Value
  Support get operations for configuration values.

  Rule: Must have valid authentication

    Scenario: Get configuration value without authentication
      Given no authentication
      And the configuration key logging.level
      When I get the configuration key
      Then I should fail with code=unauthenticated, msg=invalid or missing authorization

    Scenario: Get configuration value with invalid authentication
      Given invalid authentication
      And the configuration key logging.level
      When I get the configuration key
      Then I should fail with code=unauthenticated, msg=invalid or missing authorization

  Rule: Only existing keys can be retrieved

    Scenario: Get an existing configuration value
      Given the configuration key logging.level
      And valid authentication
      When I get the configuration key
      Then I should receive the configuration value DEBUG

    Scenario: Get a sensitive configuration value
      Given the configuration key database.password
      And valid authentication
      When I get the configuration key
      Then I should receive the configuration value ********

    Scenario: Get non-existent configuration key
      Given a non-existent configuration key
      And valid authentication
      When I get the configuration key
      Then I should fail with code=not_found, msg=value does not exist

    Scenario: Get configuration value with an empty key
      Given an empty configuration key
      And valid authentication
      When I get the configuration key
      Then I should fail with code=invalid_argument, msg=validation error: key: value is required
