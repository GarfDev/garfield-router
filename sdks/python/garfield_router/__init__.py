"""Garfield Router Python SDK.

Wraps the OpenAI Python client to automatically route requests through
Garfield Router with cost-optimised backend selection.

Usage:
    from garfield_router import GarfieldRouter

    router = GarfieldRouter("http://localhost:8050")
    response = router.chat("Summarise this document...", tier=2)
"""

from .client import GarfieldRouter, Tier

__version__ = "0.1.0"
__all__ = ["GarfieldRouter", "Tier"]
