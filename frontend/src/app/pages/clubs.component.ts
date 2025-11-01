import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../api.service';

type Tone = 'info' | 'success' | 'error';

@Component({
  selector: 'app-clubs',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './clubs.component.html'
})
export class ClubsPageComponent implements OnInit {
  clubs: any[] = [];
  clubLoading = false;
  selectedClubId: string | null = null;
  rooms: any[] = [];
  roomLoading = false;
  selectedRoomId: string | null = null;
  messages: any[] = [];
  messageLoading = false;
  newMessage = '';
  sessions: any[] = [];
  sessionLoading = false;

  newClubName = '';
  newClubDescription = '';
  newRoomName = '';
  newSessionType: 'audio' | 'video' = 'audio';
  newSessionDate = '';
  newSessionTime = '';
  newSessionDuration = 60;
  newSessionMeetingUrl = '';
  notice: { text: string; tone: Tone } | null = null;
  reportingMessageId: string | null = null;
  reportReason = '';
  reportSubmitting = false;

  constructor(private api: ApiService) {}

  ngOnInit(): void {
    this.loadClubs();
  }

  get activeClubName(): string | null {
    if (!this.selectedClubId) return null;
    const club = this.clubs.find(c => c.id === this.selectedClubId);
    return club?.name || null;
  }

  private setNotice(text: string, tone: Tone = 'info') {
    this.notice = { text, tone };
  }

  clearNotice() {
    this.notice = null;
  }

  loadClubs() {
    this.clubLoading = true;
    this.api.listClubs().subscribe({
      next: res => {
        this.clubs = res.data || [];
        this.clubLoading = false;
        if (!this.selectedClubId && this.clubs.length) {
          this.selectClub(this.clubs[0].id);
        }
      },
      error: () => {
        this.clubs = [];
        this.clubLoading = false;
        this.setNotice('We could not load your clubs just now.', 'error');
      }
    });
  }

  selectClub(clubId: string) {
    this.selectedClubId = clubId;
    this.selectedRoomId = null;
    this.loadRooms(clubId);
    this.loadSessions(clubId);
  }

  createClub() {
    if (!this.newClubName.trim()) return;
    this.clearNotice();
    this.api.createClub(this.newClubName.trim(), this.newClubDescription.trim() || undefined).subscribe({
      next: res => {
        this.newClubName = '';
        this.newClubDescription = '';
        this.loadClubs();
        if (res?.id) this.selectClub(res.id);
        this.setNotice('Club created. Invite your friends to join.', 'success');
      },
      error: () => this.setNotice('Unable to create that club right now.', 'error')
    });
  }

  joinSelectedClub() {
    if (!this.selectedClubId) return;
    this.clearNotice();
    this.api.joinClub(this.selectedClubId).subscribe({
      next: () => {
        this.refreshRooms();
        this.setNotice('You joined the club.', 'success');
      },
      error: () => this.setNotice('Unable to join that club at the moment.', 'error')
    });
  }

  loadRooms(clubId: string) {
    this.roomLoading = true;
    this.api.listRooms(clubId).subscribe({
      next: res => {
        this.rooms = res.data || [];
        this.roomLoading = false;
        if (!this.selectedRoomId && this.rooms.length) {
          this.selectRoom(this.rooms[0].id);
        }
      },
      error: () => {
        this.rooms = [];
        this.roomLoading = false;
        this.setNotice('We could not load rooms for that club.', 'error');
      }
    });
  }

  refreshRooms() {
    if (!this.selectedClubId) return;
    this.loadRooms(this.selectedClubId);
  }

  createRoom() {
    if (!this.selectedClubId || !this.newRoomName.trim()) return;
    this.clearNotice();
    this.api.createRoom(this.selectedClubId, this.newRoomName.trim()).subscribe({
      next: res => {
        this.newRoomName = '';
        this.loadRooms(this.selectedClubId!);
        if (res?.id) this.selectRoom(res.id);
        this.setNotice('Room created. Start a discussion or schedule a session.', 'success');
      },
      error: () => this.setNotice('Unable to create that room right now.', 'error')
    });
  }

  selectRoom(roomId: string) {
    this.selectedRoomId = roomId;
    this.loadMessages(roomId);
  }

  loadMessages(roomId: string) {
    this.messageLoading = true;
    this.api.getRoomMessages(roomId).subscribe({
      next: res => {
        this.messages = res.data || [];
        this.messageLoading = false;
      },
      error: () => {
        this.messages = [];
        this.messageLoading = false;
        this.setNotice('We could not load the conversation just now.', 'error');
      }
    });
  }

  refreshMessages() {
    if (!this.selectedRoomId) return;
    this.loadMessages(this.selectedRoomId);
  }

  sendMessage() {
    if (!this.selectedRoomId || !this.newMessage.trim()) return;
    this.clearNotice();
    this.api.postMessage(this.selectedRoomId, this.newMessage.trim()).subscribe({
      next: () => {
        this.newMessage = '';
        this.loadMessages(this.selectedRoomId!);
        this.setNotice('Message shared with the club.', 'success');
      },
      error: () => this.setNotice('Unable to send that message right now.', 'error')
    });
  }

  react(messageId: string, reaction: string) {
    this.api.reactToMessage(messageId, reaction).subscribe({
      next: () => this.selectedRoomId && this.loadMessages(this.selectedRoomId),
      error: () => this.setNotice('We could not update your reaction.', 'error')
    });
  }

  beginReport(messageId: string) {
    this.reportingMessageId = messageId;
    this.reportReason = '';
  }

  cancelReport() {
    this.reportingMessageId = null;
    this.reportReason = '';
    this.reportSubmitting = false;
  }

  submitReport() {
    if (!this.reportingMessageId) return;
    this.reportSubmitting = true;
    this.api.reportMessage(this.reportingMessageId, this.reportReason.trim() || undefined).subscribe({
      next: () => {
        this.setNotice('Thank you. We have logged the report for review.', 'success');
        this.cancelReport();
        if (this.selectedRoomId) this.loadMessages(this.selectedRoomId);
      },
      error: () => {
        this.setNotice('Unable to submit that report right now.', 'error');
        this.reportSubmitting = false;
      }
    });
  }

  createSession() {
    if (!this.selectedClubId) return;
    const timestamp = this.combineDateTime(this.newSessionDate, this.newSessionTime);
    if (!timestamp) {
      this.setNotice('Please provide a valid session date and time.', 'error');
      return;
    }
    this.clearNotice();
    this.api.createSession(this.selectedClubId, {
      sessionType: this.newSessionType,
      scheduledAt: timestamp,
      durationMinutes: this.newSessionDuration,
      meetingUrl: this.newSessionMeetingUrl || undefined
    }).subscribe({
      next: () => {
        this.newSessionMeetingUrl = '';
        this.loadSessions(this.selectedClubId!);
        this.setNotice('Session scheduled for the club.', 'success');
      },
      error: () => this.setNotice('Unable to schedule that session right now.', 'error')
    });
  }

  loadSessions(clubId: string) {
    this.sessionLoading = true;
    this.api.listSessions(clubId).subscribe({
      next: res => {
        this.sessions = res.data || [];
        this.sessionLoading = false;
      },
      error: () => {
        this.sessions = [];
        this.sessionLoading = false;
        this.setNotice('We could not load upcoming sessions right now.', 'error');
      }
    });
  }

  refreshSessions() {
    if (!this.selectedClubId) return;
    this.loadSessions(this.selectedClubId);
  }

  joinSession(sessionId: string) {
    this.clearNotice();
    this.api.joinSession(sessionId).subscribe({
      next: () => this.setNotice('You are on the session roster. See you there!', 'success'),
      error: () => this.setNotice('Unable to join that session right now.', 'error')
    });
  }

  private combineDateTime(date: string, time: string): string | null {
    if (!date || !time) return null;
    const iso = `${date}T${time}`;
    const dt = new Date(iso);
    if (isNaN(dt.getTime())) return null;
    return dt.toISOString();
  }

}
