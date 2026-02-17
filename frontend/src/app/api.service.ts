import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../environments/environment';

@Injectable({ providedIn: 'root' })
export class ApiService {
  constructor(private http: HttpClient) {}

  private headers(): HttpHeaders {
    let h = new HttpHeaders({ 'Content-Type': 'application/json' });
    const token = localStorage.getItem('accessToken');
    if (token) {
      h = h.set('Authorization', `Bearer ${token}`);
    }
    const devUser = localStorage.getItem('devUserId');
    if (devUser) {
      h = h.set('X-User-UUID', devUser);
    }
    return h;
  }

  // Public endpoints
  getPublicBooks(q?: string) {
    const params = q ? new HttpParams().set('q', q) : undefined;
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/public/books`, { params: params || undefined });
  }

  getPublicBook(id: string) {
    return this.http.get<{ success: boolean; data: any }>(`/svc/v1/public/books/${encodeURIComponent(id)}`);
  }

  getPublicAuthors(options: { q?: string; genre?: string } = {}) {
    let params = new HttpParams();
    if (options.q) params = params.set('q', options.q);
    if (options.genre) params = params.set('genre', options.genre);
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/public/authors`, { params });
  }

  getPublicMeetups(q?: string) {
    const params = q ? new HttpParams().set('q', q) : undefined;
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/public/meetups`, { params: params || undefined });
  }

  // Books (author/publisher)
  createBook(payload: { title: string; description?: string; authorId?: string; publisherId?: string }) {
    return this.http.post<{ success: boolean; id: string }>(`/svc/v1/books`, payload, { headers: this.headers() });
  }

  patchBook(bookId: string, payload: any) {
    return this.http.patch<{ success: boolean }>(`/svc/v1/books/${encodeURIComponent(bookId)}`, payload, { headers: this.headers() });
  }

  getBook(bookId: string) {
    return this.http.get<{ success: boolean; data: any }>(`/svc/v1/books/${encodeURIComponent(bookId)}`, { headers: this.headers() });
  }

  listMyBooks() {
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/me/books`, { headers: this.headers() });
  }

  publishBook(bookId: string, releaseDate?: string) {
    return this.http.post<{ success: boolean }>(
      `/svc/v1/books/${encodeURIComponent(bookId)}/publish`,
      releaseDate ? { releaseDate } : {},
      { headers: this.headers() }
    );
  }

  unpublishBook(bookId: string) {
    return this.http.post<{ success: boolean }>(`/svc/v1/books/${encodeURIComponent(bookId)}/unpublish`, {}, { headers: this.headers() });
  }

  listBookAssets(bookId: string) {
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/books/${encodeURIComponent(bookId)}/assets`, { headers: this.headers() });
  }

  createBookAsset(bookId: string, payload: { kind: string; source: string; url: string; checksum?: string; fileSizeBytes?: number }) {
    return this.http.post<{ success: boolean; data: any }>(`/svc/v1/books/${encodeURIComponent(bookId)}/assets`, payload, { headers: this.headers() });
  }

  deleteBookAsset(bookId: string, assetId: string) {
    return this.http.delete<{ success: boolean }>(`/svc/v1/books/${encodeURIComponent(bookId)}/assets/${encodeURIComponent(assetId)}`, { headers: this.headers() });
  }

  // Reading
  getEbookUrl(bookId: string) {
    return this.http.get<{ success: boolean; data: { url: string } }>(`/svc/v1/books/${encodeURIComponent(bookId)}/ebook-url`, { headers: this.headers() });
  }

  getPage(bookId: string, pageNo: number) {
    return this.http.get<{ success: boolean; data: any }>(`/svc/v1/books/${encodeURIComponent(bookId)}/pages/${pageNo}`, { headers: this.headers() });
  }

  listPages(bookId: string) {
    return this.http.get<{ success: boolean; data: any }>(`/svc/v1/books/${encodeURIComponent(bookId)}/pages`, { headers: this.headers() });
  }

  upsertPage(
    bookId: string,
    payload: { pageNo: number; html: string; section?: string; label?: string; audioUrl?: string; videoUrl?: string }
  ) {
    return this.http.post<{ success: boolean }>(`/svc/v1/books/${encodeURIComponent(bookId)}/pages`, payload, { headers: this.headers() });
  }

  listChapters(bookId: string) {
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/books/${encodeURIComponent(bookId)}/chapters`, { headers: this.headers() });
  }

  upsertChapter(bookId: string, payload: { number: number; title: string; pageNoStart?: number; summary?: string; audioUrl?: string }) {
    return this.http.post<{ success: boolean }>(`/svc/v1/books/${encodeURIComponent(bookId)}/chapters`, payload, { headers: this.headers() });
  }

  // Profiles
  listMyAuthors() {
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/me/authors`, { headers: this.headers() });
  }

  applyAuthor(displayName: string, bio?: string) {
    return this.http.post<{ success: boolean; data: any }>(`/svc/v1/me/authors`, { displayName, bio }, { headers: this.headers() });
  }

  listMyPublishers() {
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/me/publishers`, { headers: this.headers() });
  }

  applyPublisher(name: string, description?: string) {
    return this.http.post<{ success: boolean; data: any }>(`/svc/v1/me/publishers`, { name, description }, { headers: this.headers() });
  }

  // Collaborations
  listCollaborations() {
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/me/collaborations`, { headers: this.headers() });
  }

  createCollaboration(publisherId: string, authorId: string, notes?: string) {
    return this.http.post<{ success: boolean; data: any }>(
      `/svc/v1/publisher-author-collaborations`,
      { publisherId, authorId, notes },
      { headers: this.headers() }
    );
  }

  respondCollaboration(id: string, decision: 'accept' | 'decline', notes?: string) {
    return this.http.post<{ success: boolean }>(
      `/svc/v1/publisher-author-collaborations/${encodeURIComponent(id)}/respond`,
      { decision, notes },
      { headers: this.headers() }
    );
  }

  // Purchases / billing
  createPurchase(bookId: string, kind: 'ebook' | 'hardcopy', provider: string = 'mpesa'): Observable<{ status: number; body: any }> {
    return new Observable(observer => {
      this.http
        .post(`/svc/v1/purchases`, { bookId, kind, provider }, { headers: this.headers(), observe: 'response' })
        .subscribe({
          next: res => observer.next({ status: res.status, body: res.body }),
          error: err => observer.next({ status: err.status || 0, body: err.error || { success: false } }),
          complete: () => observer.complete()
        });
    });
  }

  getPaymentIntent(piId: number) {
    return this.http.get<{ success: boolean; data: any }>(`${environment.apiBase}/v1/billing/payment-intents/${piId}`, {
      headers: this.headers(),
      withCredentials: true
    });
  }

  confirmPurchase(purchaseId: string) {
    return this.http.post<{ success: boolean }>(`/svc/v1/purchases/${encodeURIComponent(purchaseId)}/confirm`, {}, { headers: this.headers() });
  }

  listMyPurchases() {
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/me/purchases`, { headers: this.headers() });
  }

  listEarnings() {
    return this.http.get<{ success: boolean; data: { entries: any[]; totals: Record<string, number> } }>(`/svc/v1/me/earnings`, {
      headers: this.headers()
    });
  }

  // Clubs & community
  createClub(name: string, description?: string) {
    return this.http.post<{ success: boolean; id: string }>(`/svc/v1/clubs`, { name, description }, { headers: this.headers() });
  }

  listClubs() {
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/clubs`, { headers: this.headers() });
  }

  joinClub(clubId: string) {
    return this.http.post<{ success: boolean }>(`/svc/v1/clubs/${encodeURIComponent(clubId)}/join`, {}, { headers: this.headers() });
  }

  createRoom(clubId: string, name: string) {
    return this.http.post<{ success: boolean; id: string }>(`/svc/v1/clubs/${encodeURIComponent(clubId)}/rooms`, { name }, { headers: this.headers() });
  }

  listRooms(clubId: string) {
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/clubs/${encodeURIComponent(clubId)}/rooms`, { headers: this.headers() });
  }

  getRoomMessages(roomId: string) {
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/rooms/${encodeURIComponent(roomId)}/messages`, { headers: this.headers() });
  }

  postMessage(roomId: string, content: string) {
    return this.http.post<{ success: boolean }>(`/svc/v1/rooms/${encodeURIComponent(roomId)}/messages`, { content }, { headers: this.headers() });
  }

  reactToMessage(messageId: string, reaction: string) {
    return this.http.post<{ success: boolean }>(`/svc/v1/messages/${encodeURIComponent(messageId)}/react`, { reaction }, { headers: this.headers() });
  }

  removeReaction(messageId: string, reaction: string) {
    return this.http.delete<{ success: boolean }>(`/svc/v1/messages/${encodeURIComponent(messageId)}/reactions/${encodeURIComponent(reaction)}`, {
      headers: this.headers()
    });
  }

  reportMessage(messageId: string, reason?: string) {
    return this.http.post<{ success: boolean }>(`/svc/v1/messages/${encodeURIComponent(messageId)}/report`, { reason }, { headers: this.headers() });
  }

  muteMember(clubId: string, targetUserId: string, durationMinutes?: number) {
    return this.http.post<{ success: boolean }>(
      `/svc/v1/clubs/${encodeURIComponent(clubId)}/mute`,
      { targetUserId, durationMinutes },
      { headers: this.headers() }
    );
  }

  unmuteMember(clubId: string, targetUserId: string) {
    return this.http.post<{ success: boolean }>(`/svc/v1/clubs/${encodeURIComponent(clubId)}/mute`, { targetUserId, remove: true }, { headers: this.headers() });
  }

  blockMember(clubId: string, targetUserId: string) {
    return this.http.post<{ success: boolean }>(`/svc/v1/clubs/${encodeURIComponent(clubId)}/block`, { targetUserId }, { headers: this.headers() });
  }

  unblockMember(clubId: string, targetUserId: string) {
    return this.http.post<{ success: boolean }>(`/svc/v1/clubs/${encodeURIComponent(clubId)}/block`, { targetUserId, remove: true }, { headers: this.headers() });
  }

  createSession(clubId: string, payload: { sessionType: 'audio' | 'video'; scheduledAt: string; durationMinutes?: number; bookId?: string; meetingUrl?: string; recordingUrl?: string }) {
    return this.http.post<{ success: boolean; data: any }>(`/svc/v1/clubs/${encodeURIComponent(clubId)}/sessions`, payload, { headers: this.headers() });
  }

  listSessions(clubId: string) {
    return this.http.get<{ success: boolean; data: any[] }>(`/svc/v1/clubs/${encodeURIComponent(clubId)}/sessions`, { headers: this.headers() });
  }

  joinSession(sessionId: string) {
    return this.http.post<{ success: boolean }>(`/svc/v1/sessions/${encodeURIComponent(sessionId)}/join`, {}, { headers: this.headers() });
  }

  // Reviews
  getBookReviews(bookId: string) {
    return this.http.get<{ success: boolean; data: { reviews: any[]; averageRating: number | null; totalCount: number } }>(
      `/svc/v1/books/${encodeURIComponent(bookId)}/reviews`
    );
  }

  submitReview(bookId: string, payload: { rating: number; title?: string; body?: string }) {
    return this.http.post<{ success: boolean; data: { id: string } }>(
      `/svc/v1/books/${encodeURIComponent(bookId)}/reviews`, payload, { headers: this.headers() }
    );
  }

  deleteReview(bookId: string) {
    return this.http.delete<{ success: boolean }>(
      `/svc/v1/books/${encodeURIComponent(bookId)}/reviews`, { headers: this.headers() }
    );
  }

  getMyReview(bookId: string) {
    return this.http.get<{ success: boolean; data: any }>(
      `/svc/v1/books/${encodeURIComponent(bookId)}/my-review`, { headers: this.headers() }
    );
  }
}
