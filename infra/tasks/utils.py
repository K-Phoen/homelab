from os import path

def resource_path(file):
    return path.join(path.dirname(__file__), "..", file)
