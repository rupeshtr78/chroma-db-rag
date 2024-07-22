import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()


class BinaryDistribution(setuptools.Distribution):
    def has_ext_modules(_):
        return True


setuptools.setup(
    name="embedding.py",
    version="0.1.0",
    author="gopy",
    author_email="gopy@example.com",
    description="",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/go-python/gopy",
    packages=setuptools.find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: BSD License",
        "Operating System :: OS Independent",
    ],
    include_package_data=True,
    distclass=BinaryDistribution,
)
