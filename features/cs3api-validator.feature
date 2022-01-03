Feature: list home space
  As a consumer of the cs3api
  I want to obtain my personal space
  So that I can store files

  Background: User Login
    Given user "admin" has logged in with password "admin"
    Given user "marie" has logged in with password "radioactivity"

  Scenario: First Login
    Given user "admin" has created a personal space
    When user "admin" lists all available spaces
    Then one personal space should be listed in the response

  Scenario: First Login
    Given user "marie" has created a personal space
    When user "marie" lists all available spaces
    Then one personal space should be listed in the response
