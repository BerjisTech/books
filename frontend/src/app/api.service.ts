import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../environments/environment';

@Injectable({ providedIn: 'root' })
export class ApiService {
  constructor(private http: HttpClient) {}

  private headers(): HttpHeaders {
    let h = new HttpHeaders({ 'Content-Type': 'application/json' });
    const token = localStorage.getItem('accessToken');
    if (token) { h = h.set('Authorization', `Bearer ${token}`); }
    const devUser = localStorage.getItem('devUserId');
    if (devUser) { h = h.set('X-User-ID', devUser); }
    return h;
  }

  // Public endpoints
  // Books service base for dev/prod
  private svcBase() {
    const w: any = (typeof window !== 'undefined') ? (window as any) : {};
    return w.__BOOKS_API__ || 'http://localhost:8088';
  }

  getPublicBooks(q?: string) {
    const qs = q ? `?q=${encodeURIComponent(q)}` : '';
    return this.http.get<{success: boolean; data: any[]}>(`${this.svcBase()}/v1/public/books${qs}`);
  }

  getPublicMeetups(q?: string) {
    const qs = q ? `?q=${encodeURIComponent(q)}` : '';
    return this.http.get<{success: boolean; data: any[]}>(`${this.svcBase()}/v1/public/meetups${qs}`);
  }

  // Authenticated
  createPurchase(bookId: string, kind: 'ebook'|'hardcopy', provider: string = 'mpesa'): Observable<{status: number, body: any}> {
    return new Observable(observer => {
      this.http.post(`${this.svcBase()}/v1/purchases`, { bookId, kind, provider }, { headers: this.headers(), observe: 'response' })
        .subscribe({
          next: res => observer.next({ status: res.status, body: res.body }),
          error: err => {
            const status = err.status || 0;
            observer.next({ status, body: err.error || { success: false } });
          },
          complete: () => observer.complete()
        });
    });
  }

  getPaymentIntent(piId: number) {
    return this.http.get<{success: boolean; data: any}>(`${environment.apiBase}/v1/billing/payment-intents/${piId}`, { headers: this.headers(), withCredentials: true });
  }

  confirmPurchase(purchaseId: string) {
    return this.http.post(`${this.svcBase()}/v1/purchases/${encodeURIComponent(purchaseId)}/confirm`, {}, { headers: this.headers() });
  }

  listMyPurchases() {
    return this.http.get<{success: boolean; data: any[]}>(`${this.svcBase()}/v1/me/purchases`, { headers: this.headers() });
  }
  createClub(name: string, description?: string): Observable<any> {
    return this.http.post(`/svc/v1/clubs`, { name, description }, { headers: this.headers() });
  }

  listRooms(clubId: string): Observable<any> {
    return this.http.get(`/svc/v1/clubs/${encodeURIComponent(clubId)}/rooms`, { headers: this.headers() });
  }

  postMessage(roomId: string, content: string): Observable<any> {
    return this.http.post(`/svc/v1/rooms/${encodeURIComponent(roomId)}/messages`, { content }, { headers: this.headers() });
  }
}
