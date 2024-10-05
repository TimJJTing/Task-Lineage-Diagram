# Task Lineage Diagram Interface

## Usage

To use the interface, you should:

1. Put `dot.svg` and `reachability.json` generated from the same batch into this folder. You can also use the example files from the project root, just remember to rename them to `dot.svg` and `reachability.json`.
2. Use a local http server to host these files, e.g.:

   ```sh
   python3 -m http.server 3000
   ```

3. Now you should be able to use it at localhost:3000.
