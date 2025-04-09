import { inject, Injectable } from "@angular/core";
import { WebSocketService } from "./websocket.service";
import { Observable, Subject, Subscription, takeUntil } from "rxjs";
import { HttpClient } from "@angular/common/http";
import * as uuid from 'uuid';

/**
 * Service for handling event-driven requests using WebSocket and HTTP.
 * This service makes a http request and if succesful automatically waits for a response on the websocket
 * This enabled us to use the http -> websocket pattern as if it was a normal http request.
 */
@Injectable({
  providedIn: 'root'
})
export class RequestService {
  private _websocket = inject(WebSocketService);
  private _http = inject(HttpClient);

  private handleRequest<T>(requestId: string, httpRequest: Observable<T>): Observable<T> {
    return new Observable<T>((observer) => {
      const sub = new Subscription();
      const takeUntilSub = new Subject<void>();
      sub.add(this._websocket.socketState$.pipe(takeUntil(takeUntilSub)).subscribe((state) => {
        // wait until the socket is ready
        if (state === 'READY') {
          takeUntilSub.next();
          // wait for socket response
          const takeUntilSub2 = new Subject<void>();
          sub.add(this._websocket.socketState$.pipe(takeUntil(takeUntilSub)).subscribe((state) => {
            if (state === 'CLOSED') {
              // if the socket is closed, we will never get a response, so we need to unsubscribe and throw an error
              takeUntilSub2.next();
              sub.unsubscribe();
              observer.error(new Error('WebSocket connection closed'));
              observer.complete();
            }
            sub.add(this._websocket.socket$.subscribe({
              next: (data) => {
                if (data['requestId'] === requestId) {
                  takeUntilSub2.next();
                  sub.unsubscribe();
                  observer.next(data);
                  observer.complete();
                }
              }
            }));
          }));
          httpRequest.subscribe({
            next: () => {},
            error: (err: any) => {
              observer.error(err);
              observer.complete();
              sub.unsubscribe();
            }
          });
        }
      }));
    });
  }

  private createRequest<T>(method: 'GET' | 'POST' | 'PUT' | 'DELETE', url: string, body?: any): Observable<T> {
    const requestId = uuid.v4();
    const options = {
      headers: { 'x-request-id': requestId },
      body: method === 'POST' || method === 'PUT' ? body : undefined
    };
    const httpRequest = this._http.request<T>(method, url, options);
    return this.handleRequest(requestId, httpRequest);
  }

  get<T>(url: string): Observable<T> {
    return this.createRequest('GET', url);
  }

  post<T>(url: string, body: any): Observable<T> {
    return this.createRequest('POST', url, body);
  }

  put<T>(url: string, body: any): Observable<T> {
    return this.createRequest('PUT', url, body);
  }

  delete<T>(url: string): Observable<T> {
    return this.createRequest('DELETE', url);
  }
}
