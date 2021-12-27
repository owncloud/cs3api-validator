Feature: list home space
  To store files
  As a consumer of the cs3api
  I need to be able obtain my personal space

  Scenario: Simple Login
    Given a personal space was created for user admin
    When I call listStorage Spaces
    Then there should be one personal space in the response