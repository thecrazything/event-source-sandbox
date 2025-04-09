import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { subscriptionIdInterceptor } from './interceptor/subscription-id.interceptor';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { HashLocationStrategy, LocationStrategy } from '@angular/common';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes),
    {provide: LocationStrategy, useClass: HashLocationStrategy},
    provideAnimationsAsync(),
    provideHttpClient(
      withInterceptors([
        subscriptionIdInterceptor
      ])
    )
  ]
};
