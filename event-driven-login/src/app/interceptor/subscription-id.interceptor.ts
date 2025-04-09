import { HttpEvent, HttpHandlerFn, HttpInterceptorFn, HttpRequest } from '@angular/common/http';
import { inject } from '@angular/core';
import { Observable } from 'rxjs';
import { WebSocketService } from '../service/websocket.service';

export const subscriptionIdInterceptor: HttpInterceptorFn = (
  req: HttpRequest<any>,
  next: HttpHandlerFn
): Observable<HttpEvent<any>> => {
  const socketService = inject(WebSocketService);
  const modifiedReq = req.clone({
    setHeaders: {
      'x-subscriber-id': socketService.subscriberId()
    }
  });
  return next(modifiedReq);
};
