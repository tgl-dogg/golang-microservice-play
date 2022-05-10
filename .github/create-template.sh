# read the workflow template
WORKFLOW_TEMPLATE=$(cat .github/ci-template.txt)

# iterate each MICROSERVICE in main folder
for MICROSERVICE in *-microservice/; do
    echo "generating workflow for /${MICROSERVICE}"

    # replace template MICROSERVICE placeholder with MICROSERVICE name
    WORKFLOW=$(echo "${WORKFLOW_TEMPLATE}" | sed "s/{{MICROSERVICE}}/${MICROSERVICE}/g")

    # save workflow to .github/workflows/{MICROSERVICE}
    echo "${WORKFLOW}" > .github/workflows/${MICROSERVICE}.yaml
done
