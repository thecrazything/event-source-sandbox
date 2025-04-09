import { computed, Injectable, signal } from "@angular/core";
import { BehaviorSubject, Subject } from "rxjs";

export enum WebSocketState {
  CONNECTING = 'CONNECTING',
  OPEN = 'OPEN',
  READY = 'READY',
  CLOSING = 'CLOSING',
  CLOSED = 'CLOSED'
}

@Injectable({
  providedIn: 'root'
})
export class WebSocketService {
  private socket: WebSocket | null = null;
  private reconnectInterval = 5000; // Retry every 5 seconds
  private websocketUrl = 'ws://localhost:8081/ws';
  private _subscriptionId = signal('');
  subscriberId = this._subscriptionId.asReadonly();
  private _socketSubject = new Subject<any>();
  socket$ = this._socketSubject.asObservable();
  ready = computed(() => this._subscriptionId() !== '' && this.socket !== null);
  private _socketState = new BehaviorSubject<WebSocketState>(WebSocketState.CLOSED);
  socketState$ = this._socketState.asObservable();

  connect(): void {
    if (this.socket) {
      console.warn('WebSocket is already connected.');
      return;
    }

    this.socket = new WebSocket(this.websocketUrl);

    this.socket.onopen = () => {
      console.log('WebSocket connection opened.');
      this._socketState.next(WebSocketState.OPEN);
    };

    this.socket.onmessage = (event) => {
      console.log('Message from server:', event.data);
      const json = JSON.parse(event.data);
      if (json['subscriptionId'] && Object.keys(json).length === 1) {
        this._subscriptionId.set(json['subscriptionId']);
        this._socketState.next(WebSocketState.READY);
      } else {
        this._socketSubject.next(json);
      }
    };

    this.socket.onclose = () => {
      console.log('WebSocket connection closed. Attempting to reconnect...');
      this.socket = null;
      setTimeout(() => this.connect(), this.reconnectInterval);
      this._socketState.next(WebSocketState.CLOSED);
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }
}
