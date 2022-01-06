Feature: WOPI server on public links
  As as a oCIS user without account
  I want to be able to edit office files on public links
  So that I can collaborate with regular users

  Background: Public share exists
    Given user "admin" has logged in with password "admin"
    Given user "admin" has created a resource "public-office-documents" of type "container" in the home directory with the alias "office-documents-resource"
    Given user "admin" has created a public-share "public-collaboration" with editor permissions of the resource with the alias "office-documents-resource"
    Given user "anonymous" has logged in with the token of the public-share "public-collaboration"

  Scenario: List all files in an empty public share
    When user "anonymous" lists all resources in the public-share "public-collaboration"
    Then no resource should be listed in the response

  Scenario: List all files
    Given user "anonymous" has uploaded an empty file "test.txt" to the public-share "public-collaboration"
    When user "anonymous" lists all resources in the public-share "public-collaboration"
    Then 1 resource of type "file" should be listed in the response

    Given user "anonymous" has uploaded an empty file "test2.txt" to the public-share "public-collaboration"
    When user "anonymous" lists all resources in the public-share "public-collaboration"
    Then 2 resources of type "file" should be listed in the response
