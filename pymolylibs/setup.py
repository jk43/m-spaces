from setuptools import setup, find_packages

setup(
    name='molylibs',
    packages=find_packages(include=["molylibs", "molylibs.*"]),
    version='0.1.0',
    install_requires=[
      'PyJWT',
      'pydantic',
      'fastapi',
      'grpcio',
      'pymongo',
    ],
    entry_points={
        'console_scripts': [
        ],
    },
)