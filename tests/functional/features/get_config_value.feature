Feature: Get Config Value
  Support get operations for configuration values.

  Rule: Only existing keys can be retrieved

    Scenario: Get an existing configuration value
      Given the configuration key logging.level
      When I get the configuration key
      Then I should receive the configuration value DEBUG

    Scenario: Get a sensitive configuration value
      Given the configuration key database.password
      When I get the configuration key
      Then I should receive the configuration value ********

    Scenario: Get non-existent configuration key
      Given a non-existent configuration key
      When I get the configuration key
      Then I should fail with code=not_found, msg=value does not exist

    Scenario: Get configuration value with an empty key
      Given an empty configuration key
      When I get the configuration key
      Then I should fail with code=invalid_argument, msg=validation error: key: value is required
