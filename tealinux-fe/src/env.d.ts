/// <reference types="astro/client" />

// Type definitions for Astro locals
import type { JWTPayload } from './lib/jwt';

declare namespace App {
    interface Locals {
        user?: JWTPayload;
    }
}
