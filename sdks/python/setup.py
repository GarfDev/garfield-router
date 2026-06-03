from setuptools import setup, find_packages

setup(
    name="garfield-router",
    version="0.3.0",
    description="Python SDK for Garfield Router - intelligent LLM proxy",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    author="Garfield",
    author_email="dev@garfield.co.uk",
    url="https://github.com/GarfDev/garfield-router",
    packages=find_packages(),
    python_requires=">=3.8",
    install_requires=[],  # Zero dependencies -- uses stdlib urllib
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "License :: Other/Proprietary License",
        "Programming Language :: Python :: 3",
        "Topic :: Software Development :: Libraries",
    ],
)
