import { Component, Input, OnInit } from '@angular/core';
import { Book } from '../interfaces/book';
import { BooksComponent } from "../components/books/books.component";
import { BooksService } from '../services/books.service';
import { CommonModule } from '@angular/common';
import { BookCardComponent } from "../components/book-card/book-card.component";
import { AuthService } from '../auth.service';

@Component({
  selector: 'app-home',
  standalone: true,
  templateUrl: './home.component.html',
  imports: [CommonModule, BooksComponent, BookCardComponent]
})
export class HomePageComponent implements OnInit {

  authed = false;
  accountUrl = 'http://berjis.test/account';
  userName: string | null = null;
  pageTitle: string = 'Book Recomendations';
  books: Book[] = []

  constructor(private bookService: BooksService, private auth: AuthService) { }

  ngOnInit(): void {
    // Check auth/session against Core API
    this.auth.isAuthed().subscribe(v => this.authed = !!v);
    this.auth.user().subscribe(u => this.userName = u?.name || null);
    this.auth.check();

    this.books = this.bookService.getDummyBooks(20);
  }

}
