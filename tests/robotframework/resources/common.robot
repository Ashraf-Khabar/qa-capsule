*** Settings ***
Documentation    Shared keywords and variables for SRE Capsule Robot suites.
Library          Collections
Library          OperatingSystem
Library          String

*** Variables ***
${DEFAULT_TIMEOUT}    10s

*** Keywords ***
Log Suite Context
    [Documentation]    Logs CI metadata useful when triaging failures in SRE Capsule.
    ${run_id}=    Get Environment Variable    CI_PIPELINE_ID    default=local-run
    Log    Pipeline run id: ${run_id}    console=yes
    Log    SRE Capsule upload: set SRE_CAPSULE_URL and SRE_CAPSULE_API_KEY in CI    console=yes

Should Be Valid HTTP Status
    [Arguments]    ${status_code}    ${expected}=200
    Should Be Equal As Integers    ${status_code}    ${expected}
