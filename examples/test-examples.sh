find examples -type f -name "*.nox" | while read file; do \
    echo "Executing script: $$(basename $$file)"; \
    ./nox "$$file"; \
done