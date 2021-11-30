# Teilnahme & Quelltextverwaltung

## Das Teammates Online-Portal

Die Anmeldung zum Wettbewerb findet in diesem Jahr zum ersten Mal über das neue Online-Portal [Teammates](https://teams.informaticup.de/) für Teilnehmer am informatiCup statt. Auf diesem Online-Portal könnt Ihr Euch als Teilnehmer registrieren, zu einem gemeinsamen Team einladen und dieses Team schließlich zum Wettbewerb anmelden.

### Registrierung

Erforderlich für eine Teilnahme am informatiCup ist, dass sich jedes Mitglied eines Teams mit Namen, Email-Adresse und Passwort auf dem Online-Portal [Teammates](https://teams.informaticup.de/) registriert.

Als registrierter Teilnehmer kannst Du, wenn Du ein sichtbarer Teil der informatiCup Community werden willst, für Dein Benutzerkonto ein Profil mit weiteren Information über Dich wie zum Beispiel Deinen sozialen Netzwerken oder beruflichen oder Ausbildungsstationen anlegen.

Sobald Du ein Profil angelegt hast, erscheinst Du in der öffentlichen [Liste der informatiCup-Teilnehmer](https://teams.informaticup.de/profiles).

Das Anlegen eines Profils ist dabei komplett optional und für eine Teilnahme am informatiCup nicht erforderlich. Du kannst Dein Profil auch jederzeit löschen und nur Dein Benutzerkonto behalten.

### Anmelden eines Teams

Als registrierter Benutzer kannst Du ein Team anlegen. Gib dem Team dazu einen (kreativen) Namen und gib an von welchen Hochschulen die Teammitglieder kommen. Deine Teammitglieder kannst Du über die Email-Adresse, mit der sie sich bei [Teammates](https://teams.informaticup.de/) registriert haben, einladen. Deine Einladung erscheint im [Teammates](https://teams.informaticup.de/)-Dashboard Deiner Teamkollegen. Einladungen kannst Du annehmen oder ablehnen.

Ein Team mit mindestens 2 und höchstens 4 Teammitgliedern kannst Du für den informatiCup anmelden. Beachte dazu die Anmeldefrist des aktuellen Wettbewerbs.

Die Anmeldung zum informatiCup über [Teammates](https://teams.informaticup.de/) ersetzt die bisherige Anmeldung via Email.

### Fragen

Für Fragen zu Teammates nutze bitte die [Issues des aktuellen GitHub Repositories zum Wettbewerb](https://github.com/informatiCup/informatiCup2022/issues).

### Dein SSH-Schlüssel

In [Teammates](https://teams.informaticup.de/) kannst Du zu Deinem Benutzerkonto einen öffentlichen SSH-Schlüssel eintragen. Mit diesem Schlüssel hast Du dann Zugriff auf das Git Repository Deines Teams in dem informatiCup CI System.

#### Wie sieht der öffentliche SSH-Schlüssel aus?

Der **öffentliche** Schlüssel findet sich in einer Datei mit ```.pub```-Endung. Dieser muss **ohne Zeilenumbrüche** in [Teammates](https://teams.informaticup.de/) eingetragen werden. Der Schlüssel sieht ähnlich dem folgendem Schema aus:

```
ssh-ed25519 AAAAC3NzaC1lXXXXXXXXXXXXXXX/rcj6SyU0CNdUE/w5NjoUxQDbKuwcFugyHzYhoGx5 marcus@localhost
```

folgende Beispiele sind **KEINE** gültigen öffentlichen Schlüssel:

```
SHA256:k/8OZpOFWkhOC+s+VIXBLTcJMHsUVjSuyADxhz3CLvs marcus@localhost
```
(nur die Checksumme, **kein** Schlüssel)

```
SHA256:k/8OZpOFWkhOC+s+VIXBLTcJMHsUVjSuyADxhz3CLvs
```
(nur die Checksumme, **kein** Schlüssel)

```
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktXXXXXXXXXXXXXXXXXXXXXXXXXXXXXAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDR/63I+ksXXXXXXXXXXXXXXXXXXXXXXXXoMh82IaBseQAAAJDdSEd73UhH
ewAAAAtzc2gtZWQyNTUxOQXXXXXXXXXXXXXXXXXXXXXXXXY6FMUA2yrsHBboMh82IaBseQ
AAAEXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXTdH/rcj6SyU0CNdUE/w5NjoU
xQDbKuwcFugyHzYhoGx5AAAADXRvcHJhbmdlckBJdm8=
-----END OPENSSH PRIVATE KEY-----
```
(privater Schlüssel)

```
+--[ED25519 256]--+
|  o.  oo*O+.     |
|   + o =+.*.     |
|    * * ++..     |
|   . = X.= .     |
|  . . +.S o .    |
|   o ..  = o     |
|  .  .. . B      |
|   . ..  o +     |
|    E ..   .o    |
+----[SHA256]-----+
```
(randomart ist **kein** Schlüssel)

```
ssh-ed25519 AAAAC3NzaC1l

XXXXXXXXXXXXXX/rcj6S


yU0CNdUE/w5NjoUxQDbKuwcFugyHzYhoGx5 marcus@localhost

```
(Zeilenumbrüche sind nicht erlaubt)

## Das informatiCup CI System

Jedes angemeldete Team erhält in dem neuen informatiCup CI System ("ICUPCI") ein Git Repository. Auf dieses Repository hat jedes Mitglied des Teams Zugriff mit Git via SSH, zum Beispiel mit **git clone ssh://ci@ci.informaticup.de:/repo**

Die URL ist dabei für alle Teams die gleiche, der Zugriff auf das Repository des eigenen Teams wird über die SSH-Schlüssel der Teammitglieder geregelt.

Das ICUPCI Git Repository Eures Teams steht Euch für die Dauer des laufenden Wettbewerbs zur Verfügung: als mögliches Git "origin", für die Abgabe Eurer Lösung und für laufende Softwaretests.

### Origin

Die Verwendung des ICUPCI Git Repository als "origin" ist komplett optional, steht Euch aber für den laufenden Wettbewerb natürlich gerne zur Verfügung.

Mit folgendem Befehl könnt Ihr Euren Code regelmäßig in das ICUPCI übertragen, falls Ihr das ICUPCI Git Repository nicht als "origin" nutzt:

```
git push --mirror ci@ci.informaticup.de:/repo
```

### Abgabe

Für die Einreichung Eurer Lösung am Ende des Wettbewerbs ist die Verwendung des ICUPCI Git Repository erforderlich. Ihr müsst also mindestens einmal Eure Lösung pushen, um sie einzureichen.

Beachtet das erwartete Format Eurer Lösungseinreichung in der aktuellen Aufgabenbeschreibung.

### Continuous Integration

Wenn Ihr Eure Software in das ICUPCI Git Repository Eures Teams gepusht habt, führt das informatiCup CI System automatisch Tests mit Eurer Software durch. Die Ergebnisse dieser Softwaretests könnt Ihr Euch in [Teammates](https://teams.informaticup.de/) auf der Seite Eures Teams ansehen. Ihr werdet damit also schon während der Wettbewerb noch läuft über die funktionale Korrektheit und Performance Eurer Software informiert.

## Einrichten von Git und SSH unter Windows

Der Zugriff auf das Git Repository Eures Teams im informatiCup CI System erfolgt mit Git via SSH mit dem öffentlichen SSH-Schlüssel den jedes Teammitglied für sich in [Teammates](https://teams.informaticup.de/) eintragen kann.

Wichtig für einen funktionierenden Zugriff unter Windows ist ein laufender SSH Agent, eine aktive Identity für den SSH-Schlüssel aus [Teammates](https://teams.informaticup.de/) sowie die Konfiguration von Git bzw. GitHub Desktop mit den passenden SSH Executables.

### GitHub Desktop

Für den Zugriff auf das ICUPCI Git Repository Eures Teams mit [GitHub Desktop](https://desktop.github.com/)...

1. Erzeuge mit **ssh-keygen -t ed25519 -C "your_email@example.com"** einen SSH-Schlüssel

2. Trage den öffentlichen Schlüssel in [Teammates](https://teams.informaticup.de/) ein

3. Starte den Windows Service "OpenSSH Authentication Agent". Setze den _Startup Type_ auf _Automatic_

4. Füge den privaten SSH-Schlüssel mit **ssh-add ~/.ssh/id_ed25519** (ggf. anderer Dateiname) dem OpenSSH Authentication Agent hinzu

5. In GitHub Desktop, im Menü _File > Options > Advanced_ wähle für _SSH_ die Option "Use system OpenSSH (recommended)"

6. Clone / Push mit GitHub Desktop das Repository mit der URL **ssh://ci@ci.informaticup.de:/repo**

### git

Für den Zugriff auf das ICUPCI Git Repository Eures Teams mit [git](https://git-scm.com/download/win)...

1. Erzeuge mit **ssh-keygen -t ed25519 -C "your_email@example.com"** einen SSH-Schlüssel

2. Trage den öffentlichen Schlüssel in [Teammates](https://teams.informaticup.de/) ein

3. Starte den Windows Service "OpenSSH Authentication Agent". Setze den _Startup Type_ auf _Automatic_

4. Füge den privaten SSH-Schlüssel mit **ssh-add ~/.ssh/id_ed25519** (ggf. anderer Dateiname) dem OpenSSH Authentication Agent hinzu

5. Konfiguriere git nicht die eigenen sondern die Windows SSH Executables zu verwenden. Führe dazu in der Powershell dieses Kommando aus: **[Environment]::SetEnvironmentVariable("GIT_SSH", "$((Get-Command ssh).Source)", [System.EnvironmentVariableTarget]::User)**

6. Clone / Push mit git das Repository mit der URL **ssh://ci@ci.informaticup.de:/repo**

### Fragen

Für Fragen zum informatiCup CI System nutze bitte die [Issues des aktuellen GitHub Repositories zum Wettbewerb](https://github.com/informatiCup/informatiCup2022/issues).
