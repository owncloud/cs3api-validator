Feature: list home space
  As a consumer of the cs3api
  I want to obtain my personal space
  So that I can store files

  Background: User Login
    Given user "admin" has logged in with password "admin"

  Scenario: First Login
    Given user "admin" has created a personal space with the alias "Admin Home"
    When user "admin" lists all available spaces
    Then one personal space should be listed in the response
